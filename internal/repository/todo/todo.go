package todoRepo

import (
	"database/sql"
	"fmt"
	"time"
	_ "todo/docs"
	"todo/internal/database/pg"
	"todo/internal/entity"
)

type TodoRepository struct {
	db *sql.DB
}

// @Summary Get a list of todos
// @Description Get a list of todos
// @Tags todos
// @Accept  json
// @Produce  json
// @Success 200 {array} Todo
// @Router /todos [get]
func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *entity.Todo) (int, error) {
	var itemID int
	query := fmt.Sprintf("INSERT INTO %s (user_id, title, description, completed) VALUES ($1, $2, $3, $4) RETURNING id", pg.TodosTable)
	row := r.db.QueryRow(query, todo.User_id, todo.Title, todo.Description, todo.Completed)
	err := row.Scan(&itemID)
	return itemID, err
}

func (r *TodoRepository) GetTodoById(todoId int, userId int) (*entity.TodoGetRes, error) {
	var todo entity.TodoGetRes

	query := fmt.Sprintf("SELECT id, title, description, completed FROM %s t WHERE t.id = $1 AND t.user_id = $2", pg.TodosTable)
	row := r.db.QueryRow(query, todoId, userId)
	err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
	return &todo, err
}
func (r *TodoRepository) GetTodos(userId int) ([]entity.TodoGetRes, error) {

	query := fmt.Sprintf("SELECT id, title, description, completed FROM %s t WHERE t.user_id = $1", pg.TodosTable)
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var todos []entity.TodoGetRes
	for rows.Next() {
		var todo entity.TodoGetRes
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, err
}

func (r *TodoRepository) UpdateTodo(todo *entity.TodoUpdateReq) (*entity.TodoUpdateRes, error) {
	query := fmt.Sprintf("UPDATE %s SET title=$1, description=$2, completed=$3, updated_at=$4 WHERE id = $5 AND user_id=$6", pg.TodosTable)
	res, err := r.db.Exec(query, todo.Title, todo.Description, todo.Completed, time.Now(), todo.Id, todo.UserId)

	if result, err := res.RowsAffected(); result == 0 || err != nil {
		return &entity.TodoUpdateRes{}, sql.ErrNoRows
	}

	return &entity.TodoUpdateRes{
		Id:          todo.Id,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}, err
}
func (r *TodoRepository) DeleteTodo(todoId int, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", pg.TodosTable)
	res, err := r.db.Exec(query, todoId, userId)
	if err != nil {
		return err
	}
	if result, err := res.RowsAffected(); result == 0 || err != nil {
		return sql.ErrNoRows
	}

	return nil
}
