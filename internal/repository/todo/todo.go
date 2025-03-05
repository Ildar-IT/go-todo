package todoRepo

import (
	"log/slog"
	"todo/internal/database/pg"
)

type TodoRepository struct {
	log     *slog.Logger
	storage *pg.Storage
}

func NewTodoRepository(log *slog.Logger, storage *pg.Storage) *CartRepository {
	return &TodoRepository{log: log, storage: storage}
}

func (r *TodoRepository) GetTodos(userID int) (string, error) {
	return "", nil
}
