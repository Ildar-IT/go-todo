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
	User_id     int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"status"`
}
