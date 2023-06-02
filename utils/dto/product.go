package dto

import "time"

type ProductDTO struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	Code       string    `json:"code"`
	CategoryID int       `json:"category_id"`
	Quantity   int       `json:"quantity"`
	Unit       string    `json:"unit"`
	Price      int       `json:"price"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}
