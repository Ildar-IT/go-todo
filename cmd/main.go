package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"todo/internal/config"
	"todo/internal/config/logger"
	"todo/internal/database/pg"
	"todo/internal/handler"
	jwtUtils "todo/internal/lib/jwt"
	"todo/internal/repository"
	"todo/internal/service"
)

func main() {
	//load config
	cfg := config.LoadConfig()

	//init logger
	logger := logger.SetupLogger(cfg.Env)
	logger.Info("Server started ", slog.Int("port", cfg.Http.Port))

	//connect db
	storage, err := pg.New(config.GetDbConnectionStr())
	if err != nil {
		panic(err.Error())
	}
	defer storage.DB.Close()
	jwt := jwtUtils.New(&cfg.Jwt)

	repo := repository.NewRepository(storage.DB)
	service := service.NewService(logger, repo, jwt, cfg.Salt)
	handler := handler.NewHandler(logger, service)

	router := handler.InitRoutes()
	server := http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:     router,
		ReadTimeout: cfg.Http.Timeout,
		IdleTimeout: cfg.Http.IdleTimeout,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}

}
