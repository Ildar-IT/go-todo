package handler

import (
	"log/slog"
	"net/http"
	authHandler "todo/internal/handler/auth"
	todoHandler "todo/internal/handler/todo"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/middleware"
	"todo/internal/service"
)

type Handler struct {
	todoHandler *todoHandler.TodoHandler
	authHandler *authHandler.AuthHandler
	jwt         *jwtUtils.Jwt
	log         *slog.Logger
}

func NewHandler(log *slog.Logger, services *service.Service, jwt *jwtUtils.Jwt) *Handler {
	return &Handler{
		todoHandler: todoHandler.NewTodoHandler(log, services),
		authHandler: authHandler.NewAuthHandler(log, services),
		jwt:         jwt,
		log:         log,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {

	router := http.NewServeMux()
	router.HandleFunc("POST /todo", middleware.AuthMiddleware(h.todoHandler.CreateTodo(), h.jwt, h.log))
	router.HandleFunc("GET /todo", middleware.AuthMiddleware(h.todoHandler.GetTodo(), h.jwt, h.log))
	router.HandleFunc("PATCH /todo", middleware.AuthMiddleware(h.todoHandler.UpdateTodo(), h.jwt, h.log))

	router.HandleFunc("POST /auth/login", h.authHandler.Login())
	router.HandleFunc("POST /auth/register", h.authHandler.Register())
	router.HandleFunc("POST /auth/access", middleware.RefreshTokenMiddleware(h.authHandler.UpdateAccessToken(), h.jwt, h.log))
	return router
}
