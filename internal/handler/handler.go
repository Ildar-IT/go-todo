package handler

import (
	"log/slog"
	"net/http"
	todoHandler "todo/internal/handler/todo"
	userHandler "todo/internal/handler/user"
	"todo/internal/service"
)

type Handler struct {
	todoHandler *todoHandler.TodoHandler
	userHandler *userHandler.UserHandler
}

func NewHandler(log *slog.Logger, services *service.Service) *Handler {
	return &Handler{
		todoHandler: todoHandler.NewTodoHandler(log, services),
		userHandler: userHandler.NewUserHandler(log, services),
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {

	router := http.NewServeMux()
	router.HandleFunc("POST /todo", h.todoHandler.CreateTodo())

	router.HandleFunc("POST /login", h.userHandler.Login())
	router.HandleFunc("POST /register", h.userHandler.Register())
	return router
}
