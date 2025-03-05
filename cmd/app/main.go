package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"todo/internal/app"
	"todo/internal/config"
	"todo/internal/config/logger"
	"todo/internal/database/pg"
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
	//init router
	router := http.NewServeMux()

	server := http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:     router,
		ReadTimeout: cfg.Http.Timeout,
		IdleTimeout: cfg.Http.IdleTimeout,
	}

	app.New(logger, cfg, storage, router)

	err = server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}

}
