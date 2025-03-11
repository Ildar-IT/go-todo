package roleRepo

import (
	"database/sql"
	"fmt"
	"todo/internal/database/pg"
	"todo/internal/entity"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}
func (r *RoleRepository) GetUserRole(userId int) (*entity.Role, error) {
	var role entity.Role
	query := fmt.Sprintf("SELECT r.id, r.name FROM %s ur JOIN %s r ON ur.user_id = r.id WHERE ur.user_id = $1", pg.UserRolesTable, pg.RolesTable)
	row := r.db.QueryRow(query, userId)
	err := row.Scan(&role.Id, &role.Name)
	return &role, err
}

func (r *RoleRepository) SetUserRole(userId int) (*entity.Role, error) {
	const defaultRole = "user"
	var role entity.Role
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE name = $1", pg.RolesTable)
	row := r.db.QueryRow(query, defaultRole)
	err := row.Scan(&role.Id, &role.Name)
	if err != nil {
		return &role, err
	}
	//	var res int
	query = fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES ($1, $2)", pg.UserRolesTable)
	row = r.db.QueryRow(query, userId, role.Id)
	err = row.Scan()
	if err != nil {
		return &role, err
	}
	return &role, err
}

// CREATE TABLE user_roles (
//     user_id INT REFERENCES users(id) ON DELETE CASCADE,
//     role_id INT REFERENCES roles(id) ON DELETE CASCADE,
//     PRIMARY KEY (user_id, role_id)
// );

// CREATE TABLE roles(
//     id SERIAL PRIMARY KEY,
//     name VARCHAR(50) UNIQUE NOT NULL
// );
