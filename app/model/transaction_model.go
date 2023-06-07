package model

import "time"

type Transaction struct {
	ID         int       `json:"id"`
	OrderID    int       `json:"order_id"`
	Status     string    `json:"status"`
	Payment    string    `json:"payment"`
	MemberCode string    `json:"member_code"`
	Amount     int       `json:"amount"`
	Service    int       `json:"service"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User      `json:"user"`
}
