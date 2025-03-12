package service

import (
	"log/slog"
	"todo/internal/entity"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/repository"
	todoService "todo/internal/service/todo"
	userService "todo/internal/service/user"
)

type Todo interface {
	CreateTodo(todo *entity.Todo) (int, error)
}

type User interface {
	Login(user *entity.UserLoginReq) (entity.TokensRes, error, int)
	Register(user *entity.UserRegisterReq) (entity.TokensRes, error, int)
	GenerateTokens(userId int, role string) (entity.TokensRes, error, int)
}

type Service struct {
	Todo
	User
}

func NewService(log *slog.Logger, repo *repository.Repository, jwt *jwtUtils.Jwt, salt string) *Service {
	return &Service{
		Todo: todoService.NewTodoService(log, repo.Todo),
		User: userService.NewUserService(log, repo.User, repo.Role, jwt, salt),
	}
}
