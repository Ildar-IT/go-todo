package repository

import (
	"database/sql"
	"todo/internal/entity"
	roleRepo "todo/internal/repository/role"
	todoRepo "todo/internal/repository/todo"
	userRepo "todo/internal/repository/user"
)

type Todo interface {
	Create(todo *entity.Todo) (int, error)
	GetTodoById(todoId int, userId int) (*entity.TodoGetRes, error)
	UpdateTodo(todo *entity.TodoUpdateReq) (*entity.TodoUpdateRes, error)
	GetTodos(userId int) ([]entity.TodoGetRes, error)
	DeleteTodo(todoId int, userId int) error
}

type User interface {
	Create(user *entity.User) (int, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUsers() ([]entity.User, error)
	GetUsersWithTasksByDate(date string) (map[int]*entity.UserTasks, error)
}

type Role interface {
	GetUserRole(userId int) (*entity.Role, error)
	SetUserRole(userId int) (*entity.Role, error)
}

type Repository struct {
	Todo
	User
	Role
}

func NewRepository(db *sql.DB) *Repository {

	return &Repository{
		Todo: todoRepo.NewTodoRepository(db),
		User: userRepo.NewUserRepository(db),
		Role: roleRepo.NewRoleRepository(db),
	}
}
