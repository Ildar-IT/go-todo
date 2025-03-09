package todoRepo

import (
	"database/sql"
	"fmt"
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
	createItemQuery := fmt.Sprintf("INSERT INTO %s (user_id, title, description, completed) VALUES ($1, $2, $3, $4) RETURNING id", pg.TodosTable)
	row := r.db.QueryRow(createItemQuery, todo.User_id, todo.Title, todo.Description, todo.Completed)
	err := row.Scan(&itemID)
	return itemID, err
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
