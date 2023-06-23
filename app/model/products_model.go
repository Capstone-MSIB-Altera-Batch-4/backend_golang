package model

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	ProductID   string    `json:"product_id"`
	CategoryID  int       `json:"category_id"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
