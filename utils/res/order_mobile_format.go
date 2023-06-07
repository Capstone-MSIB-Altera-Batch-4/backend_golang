package res

import (
	"point-of-sale/app/model"
	"strings"
)

type SetCashierOrderResponse struct {
	ID          int                     `json:"order_id"`
	OrderCode   string                  `json:"order_code"`
	Name        string                  `json:"name"`
	OrderOption string                  `json:"order_option"`
	NumberTable int                     `json:"number_table"`
	Service     int                     `json:"service"`
	Subtotal    int                     `json:"subtotal"`
	GrandTotal  int                     `json:"grand_total"`
	Items       []SetItemOutputResponse `json:"items"`
	Transaction SetTransactionResponse  `json:"transaction"`
}

type SetItemOutputResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Note     string `json:"note"`
	Subtotal int    `json:"subtotal"`
}

type SetTransactionResponse struct {
	PaymentStatus string `json:"payment_status"`
	PaymentMethod string `json:"payment_method"`
}

func TransformOrderResponse(order model.Order) SetCashierOrderResponse {
	setItems := make([]SetItemOutputResponse, len(order.Items))
	subtotal := 0

	for i, item := range order.Items {
		setItems[i] = SetItemOutputResponse{
			ID:       item.ProductID,
			Name:     item.Products.Name,
			Quantity: item.Quantity,
			Note:     item.Note,
			Price:    item.Products.Price,
			Subtotal: item.Subtotal,
		}

		subtotal += item.Subtotal
	}

	setTransaction := SetTransactionResponse{
		PaymentStatus: order.Transaction.Status,
		PaymentMethod: order.Transaction.Payment,
	}

	setResponse := SetCashierOrderResponse{
		ID:          order.ID,
		OrderCode:   order.OrderCode,
		Name:        order.Name,
		OrderOption: strings.Title(order.OrderOption),
		NumberTable: order.NumberTable,
		Service:     order.Transaction.Service,
		Subtotal:    subtotal,
		GrandTotal:  order.Transaction.Amount,
		Items:       setItems,
		Transaction: setTransaction,
	}

	return setResponse
}

// SetSearchOrderResponse Pages Order employee
type SetSearchOrderResponse struct {
	CategoryName string               `json:"category_name"`
	Products     []SetGetItemResponse `json:"products,omitempty"`
}

type SetGetItemResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Status      bool   `json:"status"`
}

func TransformItemOrder(search []model.Product) []SetGetItemResponse {
	setItems := make([]SetGetItemResponse, len(search))

	for i, product := range search {
		status := product.Quantity > 0
		setItems[i] = SetGetItemResponse{
			ID:          product.ID,
			Name:        product.Name,
			ImageURL:    product.Image,
			Description: product.Description,
			Price:       product.Price,
			Status:      status,
		}
	}

	return setItems
}

func TransformCategoryOrder(search model.Category) SetSearchOrderResponse {
	setItems := make([]SetGetItemResponse, len(search.Products))

	for i, product := range search.Products {
		status := product.Quantity > 0
		setItems[i] = SetGetItemResponse{
			ID:          product.ID,
			Name:        product.Name,
			ImageURL:    product.Image,
			Description: product.Description,
			Price:       product.Price,
			Status:      status,
		}
	}

	setResponse := SetSearchOrderResponse{
		CategoryName: search.Name,
		Products:     setItems,
	}

	return setResponse
}
