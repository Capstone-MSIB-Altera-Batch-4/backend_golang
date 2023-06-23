package model

import "time"

type User struct {
	ID           int           `json:"id"`
	UserCode     string        `json:"id_code"`
	Username     string        `json:"name"`
	Password     string        `json:"password"`
	Role         string        `json:"role"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Transactions []Transaction `json:"transactions"`
}
