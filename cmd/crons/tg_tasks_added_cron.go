package main

import (
	"todo/internal/config"
	"todo/internal/config/logger"
	"todo/internal/cron"
	"todo/internal/database/pg"
	"todo/internal/repository"
)

func main() {

	//load config
	cfg := config.LoadConfig()

	//init logger
	logger := logger.SetupLogger(cfg.Env)
	storage, err := pg.New(config.GetDbConnectionStr(cfg.EnvFilePath))
	if err != nil {
		panic(err.Error())
	}
	defer storage.DB.Close()

	repo := repository.NewRepository(storage.DB)

	cron := cron.New(logger, repo, cfg.Tg.Token, cfg.Tg.ChatId)
	cron.Run()
}
