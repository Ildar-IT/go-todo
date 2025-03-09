package main

import (
	"embed"
	"flag"
	"fmt"
	"todo/internal/config"
	"todo/internal/database/pg"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	// setup database
	storage, err := pg.New(config.GetDbConnectionStr())
	if err != nil {
		panic(err.Error())
	}
	err = storage.DB.Ping()
	if err != nil {
		panic(err.Error())
	}
	defer storage.DB.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	var migrationType string
	flag.StringVar(&migrationType, "type", "", "up down")
	flag.Parse()
	fmt.Println(migrationType)

	if migrationType == "down" {
		if err := goose.Down(storage.DB, "migrations"); err != nil {
			panic(err)
		}

		return
	}

	if err := goose.Up(storage.DB, "migrations"); err != nil {
		panic(err)
	}

}
