package dto

import "point-of-sale/app/model"

type CreateOrderRequest struct {
	Name        string               `json:"name" validate:"required"`
	OrderOption string               `json:"order_option" validate:"required"`
	TableNumber int                  `json:"number_table" validate:"numeric"`
	MemberCode  string               `json:"member_code"`
	Payment     string               `json:"payment" validate:"required"`
	Items       []CreateItemsRequest `json:"items" validate:"required"`
	User        model.User
}

type CreateItemsRequest struct {
	ProductID int    `json:"product_id" validate:"required,numeric"`
	Note      string `json:"note,omitempty"`
	Quantity  int    `json:"quantity" validate:"required,numeric"`
}

type SearchCategoryRequest struct {
	Name string `json:"name"`
}

type SearchRequest struct {
	Keyword string `json:"keyword"`
}

type CreateTokenCCRequest struct {
}
