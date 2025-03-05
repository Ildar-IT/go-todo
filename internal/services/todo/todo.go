package todoService

import (
	"log/slog"
	repository "todo/internal/database/pg/repository/todo"
)

type TodoService struct {
	log  *slog.Logger
	repo *repository.TodoRepository
}

func NewCartService(log *slog.Logger, repo *repository.TodoRepository) *TodoService {

	return &TodoService{
		log:  log,
		repo: repo,
	}
}

func (s *TodoService) GetCart(userID int) (string, error) {
	return s.repo.GetTodos(userID)
}
