package entity

import "time"

// Todo представляет собой модель задачи
// swagger:model Todo
type Todo struct {
	Id          int
	User_id     int
	Title       string
	Description string
	Completed   bool
	Created_at  time.Time
	Updated_at  time.Time
}

// TodoCreateReq представляет запрос на создание задачи
// swagger:model TodoCreateReq
type TodoCreateReq struct {
	// Название задачи
	// required: true
	// example: Купить молоко
	Title string `json:"title" validate:"required,min=6,max=25"`
	// Описание задачи
	// required: true
	// example: Купить молоко в магазине на углу
	Description string `json:"description" validate:"required,max=360"`
	// Статус выполнения задачи
	// required: false
	// example: false
	Completed bool `json:"completed"`
}

// TodoCreateRes представляет ответ на запрос создания задачи
// swagger:model TodoCreateRes
type TodoCreateRes struct {
	Id int `json:"id"`
}

// TodoGetRes представляет ответ на запрос получения задачи
// swagger:model TodoGetRes
type TodoGetRes struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// TodoUpdateReq представляет запрос на обновление задачи
// swagger:model TodoUpdateReq
type TodoUpdateReq struct {
	UserId      int    `json:"user_id" validate:"required"`
	Id          int    `json:"id"  validate:"required"`
	Title       string `json:"title" validate:"required,min=6,max=25"`
	Description string `json:"description" validate:"required,max=360"`
	Completed   bool   `json:"completed"`
}

// TodoUpdateRes представляет ответ на запрос обновления задачи
// swagger:model TodoUpdateRes
type TodoUpdateRes struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
