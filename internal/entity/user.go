package entity

import "time"

// User представляет собой модель пользователя
// swagger:model User
type User struct {
	Id            int
	Username      string
	Email         string
	Password_hash string
	Created_at    time.Time
	Updated_at    time.Time
}
type UserTasks struct {
	User  User
	Todos []Todo
}

// UserLoginReq представляет запрос на вход пользователя
// swagger:model UserLoginReq
type UserLoginReq struct {
	// Email пользователя
	// required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`
	// Пароль пользователя
	// required: true
	// example: password123
	Password string `json:"password" validate:"required,min=6"`
}

// UserRegisterReq представляет запрос на регистрацию пользователя
// swagger:model UserRegisterReq
type UserRegisterReq struct {
	// Имя пользователя
	// required: true
	// example: JohnDoe
	Username string `json:"name" validate:"required,max=32,min=2"`
	// Email пользователя
	// required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`
	// Пароль пользователя
	// required: true
	// example: password123
	Password string `json:"password" validate:"required,min=6"`
}

// TokensRes представляет ответ с токенами доступа и обновления
// swagger:model TokensRes
type TokensRes struct {
	// Токен доступа
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Access string `json:"access"`
	// Токен обновления
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Refresh string `json:"refresh"`
}

// TokenAccessRes представляет ответ с токеном доступа
// swagger:model TokenAccessRes
type TokenAccessRes struct {
	// Токен доступа
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Access string `json:"access"`
}
