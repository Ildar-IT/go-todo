package todoRepo

import (
	"database/sql"
	"fmt"
	"time"
	"todo/internal/database/pg"
	"todo/internal/entity"
)

type TodoRepository struct {
	db *sql.DB
}

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

func (r *TodoRepository) UpdateTodo(todo *entity.TodoUpdateReq) (*entity.TodoUpdateRes, error) {
	query := fmt.Sprintf("UPDATE %s SET title=$1, description=$2, completed=$3, updated_at=$4 WHERE id = $5 AND user_id=$6", pg.TodosTable)
	_, err := r.db.Exec(query, todo.Title, todo.Description, todo.Completed, time.Now(), todo.Id, todo.UserId)
	return &entity.TodoUpdateRes{
		Id:          todo.Id,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}, err
}

// CREATE TABLE tasks (
//     id SERIAL PRIMARY KEY,
//     user_id INT REFERENCES users(id) ON DELETE CASCADE,
//     title VARCHAR(255) NOT NULL,
//     description TEXT,
//     completed BOOLEAN DEFAULT FALSE,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
