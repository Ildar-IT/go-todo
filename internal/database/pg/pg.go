package pg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

const (
	TodosTable = "tasks"
	UsersTable = "users"
)

func New(connecStr string) (*Storage, error) {
	//const op = "database.pg.New"

	db, err := sql.Open("postgres", connecStr)
	if err != nil {
		return &Storage{}, err
	}

	err = db.Ping()
	if err != nil {
		return &Storage{}, err
	}

	return &Storage{DB: db}, nil
}
