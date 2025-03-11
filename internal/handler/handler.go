package handler

import (
	"log/slog"
	"net/http"
	todoHandler "todo/internal/handler/todo"
	userHandler "todo/internal/handler/user"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/middleware"
	"todo/internal/service"
)

type Handler struct {
	todoHandler *todoHandler.TodoHandler
	userHandler *userHandler.UserHandler
	jwt         *jwtUtils.Jwt
	log         *slog.Logger
}

func NewHandler(log *slog.Logger, services *service.Service, jwt *jwtUtils.Jwt) *Handler {
	return &Handler{
		todoHandler: todoHandler.NewTodoHandler(log, services),
		userHandler: userHandler.NewUserHandler(log, services),
		jwt:         jwt,
		log:         log,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {

	router := http.NewServeMux()
	router.HandleFunc("POST /todo", h.todoHandler.CreateTodo())

	router.HandleFunc("POST /login", h.userHandler.Login())
	router.HandleFunc("POST /register", middleware.AuthMiddleware(h.userHandler.Register(), h.jwt, h.log))
	return router
}
