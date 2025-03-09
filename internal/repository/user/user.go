package userRepo

import (
	"database/sql"
	"fmt"
	"todo/internal/database/pg"
	"todo/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entity.User) (int, error) {
	var userId int
	query := fmt.Sprintf("INSERT INTO %s (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id", pg.UsersTable)
	row := r.db.QueryRow(query, user.Username, user.Email, user.Password_hash)
	err := row.Scan(&userId)
	return userId, err
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var result entity.User
	query := fmt.Sprintf("SELECT id, username, email, password_hash, created_at, updated_at FROM %s WHERE email = $1", pg.UsersTable)
	row := r.db.QueryRow(query, email)
	err := row.Scan(&result.Id, &result.Username, &result.Email, &result.Password_hash, &result.Created_at, &result.Updated_at)
	return &result, err
}
