package app

import (
	"log/slog"
	"net/http"
	"todo/internal/config"
	"todo/internal/database/pg"
	todoHandler "todo/internal/handlers/todo"
	todoRepo "todo/internal/repository/todo"
	todoService "todo/internal/services/todo"
)

func New(log *slog.Logger, cfg *config.Config, storage *pg.Storage, router *http.ServeMux) {

	todoRepo := todoRepo.NewTodoRepository(log, storage)
	todoService := todoService.NewCartService(log, todoRepo)
	todoHandler.NewCartHandler(log, router, todoService)
}
