package todoService

import (
	"log/slog"
	"todo/internal/entity"
	"todo/internal/repository"
)

type TodoService struct {
	log  *slog.Logger
	repo repository.Todo
}

func NewTodoService(log *slog.Logger, repo repository.Todo) *TodoService {

	return &TodoService{
		log:  log,
		repo: repo,
	}
}

func (s *TodoService) CreateTodo(todo *entity.Todo) (int, error) {
	return s.repo.Create(todo)
}
