package model

import "time"

type User struct {
	ID        int       `json:"id"`
	UsersCode string    `json:"users_code"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
