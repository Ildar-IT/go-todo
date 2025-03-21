package handler

import (
	"log/slog"
	"net/http"
	authHandler "todo/internal/handler/auth"
	todoHandler "todo/internal/handler/todo"
	userHandler "todo/internal/handler/user"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/middleware"
	"todo/internal/service"

	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	todoHandler *todoHandler.TodoHandler
	authHandler *authHandler.AuthHandler
	userHandler *userHandler.UserHandler
	jwt         *jwtUtils.Jwt
	log         *slog.Logger
}

func NewHandler(log *slog.Logger, services *service.Service, jwt *jwtUtils.Jwt) *Handler {

	validator := validator.New()
	return &Handler{
		todoHandler: todoHandler.NewTodoHandler(log, services, validator),
		authHandler: authHandler.NewAuthHandler(log, services, validator),
		userHandler: userHandler.NewUserHandler(log, services),
		jwt:         jwt,
		log:         log,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {

	router := http.NewServeMux()
	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	authMiddleware := middleware.ChainMiddleware(
		middleware.AuthMiddleware(h.jwt, h.log),
	)
	authAndRoleMw := middleware.ChainMiddleware(
		middleware.AuthMiddleware(h.jwt, h.log),
		middleware.RoleCheckMiddleware(h.log),
	)
	router.HandleFunc("POST /todo", authMiddleware(h.todoHandler.CreateTodo()))
	router.HandleFunc("GET /todo/{id}", authMiddleware(h.todoHandler.GetTodo()))
	router.HandleFunc("PATCH /todo", authMiddleware(h.todoHandler.UpdateTodo()))
	router.HandleFunc("DELETE /todo/{id}", authMiddleware(h.todoHandler.DeleteTodo()))

	router.HandleFunc("GET /todos", authMiddleware(h.todoHandler.GetTodos()))

	router.HandleFunc("POST /auth/login", h.authHandler.Login())

	router.HandleFunc("POST /auth/register", h.authHandler.Register())
	router.HandleFunc("POST /auth/access", middleware.ChainMiddleware(middleware.RefreshTokenMiddleware(h.jwt, h.log))(h.authHandler.UpdateAccessToken()))

	router.HandleFunc("GET /users", authAndRoleMw(h.userHandler.GetAllUsers()))
	return router
}
