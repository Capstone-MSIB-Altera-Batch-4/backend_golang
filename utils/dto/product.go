package dto

type CreateProductRequest struct {
	ProductID     string `form:"products_id" validate:"required"`
	Name          string `form:"products_name" validate:"required"`
	CategoryID    int    `form:"products_category" validate:"required,numeric"`
	ProductsImage string `form:"-" validate:"required,image"`
	Quantity      int    `form:"products_quantity" validate:"required,numeric"`
	Price         int    `form:"products_price" validate:"required,numeric"`
	Unit          string `form:"products_unit" validate:"required"`
	Description   string `form:"products_description" validate:"required"`
}

type UpdateProductRequest struct {
	ProductID     string `form:"products_id" validate:"required"`
	Name          string `form:"products_name" validate:"required"`
	CategoryID    int    `form:"products_category" validate:"required,numeric"`
	ProductsImage string `form:"-" validate:"image"`
	Quantity      int    `form:"products_quantity" validate:"required,numeric"`
	Price         int    `form:"products_price" validate:"required,numeric"`
	Unit          string `form:"products_unit" validate:"required"`
	Description   string `form:"products_description" validate:"required"`
}
