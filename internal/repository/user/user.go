package userRepo

import (
	"database/sql"
	"fmt"
	"todo/internal/database/pg"
	"todo/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entity.User) (int, error) {
	var userId int
	query := fmt.Sprintf("INSERT INTO %s (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id", pg.UsersTable)
	row := r.db.QueryRow(query, user.Username, user.Email, user.Password_hash)
	err := row.Scan(&userId)
	return userId, err
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var result entity.User
	query := fmt.Sprintf("SELECT id, username, email, password_hash, created_at, updated_at FROM %s WHERE email = $1", pg.UsersTable)
	row := r.db.QueryRow(query, email)
	err := row.Scan(&result.Id, &result.Username, &result.Email, &result.Password_hash, &result.Created_at, &result.Updated_at)

	return &result, err
}

func (r *UserRepository) GetUsers() ([]entity.User, error) {
	var users []entity.User
	query := fmt.Sprintf("SELECT id, username, email, password_hash, created_at, updated_at FROM %s", pg.UsersTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password_hash, &user.Created_at, &user.Updated_at); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, err
}
func (r *UserRepository) GetUsersWithTasksByDate(date string) (map[int]*entity.UserTasks, error) {

	query := fmt.Sprintf(
		`SELECT u.id, u.username, u.email, t.id, t.title, t.description, t.completed FROM %s u 
		INNER JOIN %s t ON u.id = t.user_id AND date(t.created_at) = date($1)
		`,
		pg.UsersTable, pg.TodosTable)

	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	usersWithTodos := make(map[int]*entity.UserTasks)
	for rows.Next() {
		var user entity.User
		var todo entity.Todo

		err := rows.Scan(&user.Id, &user.Username, &user.Email, &todo.Id, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return nil, err
		}

		if _, exists := usersWithTodos[user.Id]; !exists {
			usersWithTodos[user.Id] = &entity.UserTasks{
				User: entity.User{
					Id:       user.Id,
					Username: user.Username,
					Email:    user.Email,
				},
				Todos: []entity.Todo{},
			}
		}

		usersWithTodos[user.Id].Todos = append(usersWithTodos[user.Id].Todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usersWithTodos, err
}

// func (r *TodoRepository) GetTodos(userId int) ([]entity.TodoGetRes, error) {

// 	query := fmt.Sprintf("SELECT id, title, description, completed FROM %s t WHERE t.user_id = $1", pg.TodosTable)
// 	rows, err := r.db.Query(query, userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var todos []entity.TodoGetRes
// 	for rows.Next() {
// 		var todo entity.TodoGetRes
// 		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed); err != nil {
// 			return nil, err
// 		}
// 		todos = append(todos, todo)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return todos, err
// }
