package res

import (
	"point-of-sale/app/model"
	"time"
)

type SetOrderResponse struct {
	OrderID     string                 `json:"order_id"`
	OrderOption string                 `json:"order_option"`
	Name        string                 `json:"name,omitempty"`
	NumberTable int                    `json:"number_table,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	Status      string                 `json:"status,omitempty"`
	Payment     string                 `json:"payment"`
	Subtotal    int                    `json:"sub_total,omitempty"`
	Service     int                    `json:"service,omitempty"`
	GrandTotal  int                    `json:"grand_total"`
	Items       []SetItemOrderResponse `json:"items,omitempty"`
}

type SetItemOrderResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Subtotal int    `json:"subtotal"`
}

func TransformResponse(order model.Order) SetOrderResponse {
	setItems := make([]SetItemOrderResponse, len(order.Items))
	Subtotal := 0

	for i, item := range order.Items {
		setItems[i] = SetItemOrderResponse{
			ID:       item.ID,
			Name:     item.ProductName,
			Quantity: item.Quantity,
			Price:    item.Subtotal / item.Quantity,
			Subtotal: item.Subtotal,
		}
		Subtotal += item.Subtotal

	}

	setResponse := SetOrderResponse{
		OrderID:     order.OrderCode,
		OrderOption: order.OrderOption,
		Name:        order.Name,
		NumberTable: order.NumberTable,
		CreatedAt:   order.CreatedAt,
		Payment:     order.Transaction.Payment,
		Subtotal:    Subtotal,
		Service:     order.Transaction.Service,
		GrandTotal:  order.Transaction.Amount,
		Items:       setItems,
	}

	return setResponse
}

func TransformResponseDataOrder(order model.Order) SetOrderResponse {
	setItems := make([]SetItemOrderResponse, len(order.Items))

	for i, item := range order.Items {
		setItems[i] = SetItemOrderResponse{
			ID:       item.ID,
			Name:     item.ProductName,
			Quantity: item.Quantity,
			Price:    item.Subtotal / item.Quantity,
			Subtotal: item.Subtotal,
		}

	}

	setResponse := SetOrderResponse{
		OrderID:     order.OrderCode,
		OrderOption: order.OrderOption,
		CreatedAt:   order.CreatedAt,
		Status:      order.Transaction.Status,
		Payment:     order.Transaction.Payment,
		GrandTotal:  order.Transaction.Amount,
		Items:       setItems,
	}

	return setResponse
}
