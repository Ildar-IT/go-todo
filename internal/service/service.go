package service

import (
	"log/slog"
	"todo/internal/entity"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/repository"
	authService "todo/internal/service/auth"
	todoService "todo/internal/service/todo"
	userService "todo/internal/service/user"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Todo interface {
	CreateTodo(todo *entity.Todo) (int, error, int)
	GetTodo(todoId int, userId int) (*entity.TodoGetRes, error, int)
	UpdateTodo(todo *entity.TodoUpdateReq) (*entity.TodoUpdateRes, error, int)
	GetTodos(userId int) ([]entity.TodoGetRes, error, int)
	DeleteTodo(todoId int, userId int) (error, int)
}

type Auth interface {
	Login(user *entity.UserLoginReq) (entity.TokensRes, error, int)
	Register(user *entity.UserRegisterReq) (entity.TokensRes, error, int)
	GenerateTokens(userId int, role string) (entity.TokensRes, error, int)
	GenerateAccessToken(userId int, role string) (string, error, int)
	GenerateRefreshToken(userId int, role string) (string, error, int)
}

type User interface {
	GetAllUsers() ([]entity.User, error, int)
}

type Service struct {
	Todo
	Auth
	User
}

func NewService(log *slog.Logger, repo *repository.Repository, jwt *jwtUtils.Jwt, salt string) *Service {
	return &Service{
		Todo: todoService.NewTodoService(log, repo.Todo),
		Auth: authService.NewAuthService(log, repo.User, repo.Role, jwt, salt),
		User: userService.NewUserService(log, repo.User),
	}
}
