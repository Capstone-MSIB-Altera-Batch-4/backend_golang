package model

import "time"

type OrderItems struct {
	ID          int       `json:"id"`
	OrderID     int       `json:"order_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Subtotal    int       `json:"subtotal"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
