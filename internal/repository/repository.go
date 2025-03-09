package repository

import (
	"database/sql"
	"todo/internal/entity"
	todoRepo "todo/internal/repository/todo"
	userRepo "todo/internal/repository/user"
)

type Todo interface {
	Create(todo *entity.Todo) (int, error)
}

type User interface {
	Create(user *entity.User) (int, error)
	GetUserByEmail(email string) (*entity.User, error)
}

type Repository struct {
	Todo
	User
}

func NewRepository(db *sql.DB) *Repository {

	return &Repository{
		Todo: todoRepo.NewTodoRepository(db),
		User: userRepo.NewUserRepository(db),
	}
}
