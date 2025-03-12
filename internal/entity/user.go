package entity

import "time"

type User struct {
	Id            int
	Username      string
	Email         string
	Password_hash string
	Created_at    time.Time
	Updated_at    time.Time
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokensRes struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
