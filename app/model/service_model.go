package model

import "time"

type Service struct {
	ID        int       `json:"id"`
	Service   int       `json:"service"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
