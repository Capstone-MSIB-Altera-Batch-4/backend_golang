package model

import "time"

type Order struct {
	ID          int          `json:"id"`
	OrderCode   string       `json:"order_code"`
	Name        string       `json:"name"`
	OrderOption string       `json:"order_option"`
	NumberTable int          `json:"number_table"`
	Items       []OrderItems `json:"items" gorm:"foreignKey:OrderID"`
	Transaction Transaction  `json:"transaction"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
