package model

import "time"

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Products  []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string {
	return "category"
}
