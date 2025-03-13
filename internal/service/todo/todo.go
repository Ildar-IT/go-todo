package todoService

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
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

func (s *TodoService) CreateTodo(todo *entity.Todo) (int, error, int) {
	const op = "services.todo.CreateTodo"
	log := s.log.With(
		slog.String("op", op),
	)

	id, err := s.repo.Create(todo)
	if err != nil {
		log.Error("Create Todo error", "error", err.Error())
		return 0, errors.New("create todo error"), http.StatusInternalServerError
	}
	return id, nil, http.StatusOK
}

func (s *TodoService) GetTodo(todoId int, userId int) (*entity.TodoGetRes, error, int) {
	const op = "services.todo.GetTodo"
	log := s.log.With(
		slog.String("op", op),
	)

	todo, err := s.repo.GetTodoById(todoId, userId)
	if err != nil {
		log.Error("Get todo by id error", "error", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return &entity.TodoGetRes{}, errors.New("task not found"), http.StatusBadRequest
		}
		return &entity.TodoGetRes{}, errors.New("get todo error"), http.StatusInternalServerError
	}
	return todo, nil, http.StatusOK
}

func (s *TodoService) UpdateTodo(todo *entity.TodoUpdateReq) (*entity.TodoUpdateRes, error, int) {
	const op = "services.todo.UpdateTodo"
	log := s.log.With(
		slog.String("op", op),
	)

	todoRes, err := s.repo.UpdateTodo(todo)
	if err != nil {
		log.Error("Update todo error", "error", err.Error())
		return todoRes, errors.New("update todo error"), http.StatusInternalServerError
	}
	return todoRes, nil, http.StatusOK
}
