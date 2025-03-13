package entity

import "time"

type Todo struct {
	Id          int
	User_id     int
	Title       string
	Description string
	Completed   bool
	Created_at  time.Time
	Updated_at  time.Time
}

type TodoCreateReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
type TodoCreateRes struct {
	Id int `json:"id"`
}

type TodoGetRes struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
type TodoUpdateReq struct {
	UserId      int    `json:"user_id"`
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
type TodoUpdateRes struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
