package res

import "point-of-sale/app/model"

type SetPaymentHistory struct {
	ID          int    `json:"id"`
	OrderID     string `json:"order_id"`
	Name        string `json:"name"`
	TableNumber int    `json:"table_number"`
	GrandTotal  int    `json:"grand_total"`
	Payment     string `json:"payment"`
	Status      string `json:"status"`
}

type SetDetailPaymentHistory struct {
	ID          int                           `json:"id"`
	OrderID     string                        `json:"order_id"`
	Name        string                        `json:"name"`
	TableNumber int                           `json:"table_number"`
	Item        []SetItemDetailPaymentHistory `json:"item"`
	SubTotal    int                           `json:"sub_total"`
	Service     int                           `json:"service"`
	GrandTotal  int                           `json:"grand_total"`
	Payment     string                        `json:"payment"`
	Status      string                        `json:"status"`
}

type SetItemDetailPaymentHistory struct {
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
}

func SetResponseOrderHistory(order []model.Order) []SetPaymentHistory {
	setHistoryOrder := make([]SetPaymentHistory, len(order))

	for i, m := range order {
		st := m.Transaction.Status
		if st == "PAID" {
			st = "Payment Approved"
		} else {
			st = "Payment Not Approved"
		}
		setHistoryOrder[i] = SetPaymentHistory{
			ID:          m.ID,
			OrderID:     m.OrderCode,
			Name:        m.Name,
			TableNumber: m.NumberTable,
			GrandTotal:  m.Transaction.Amount,
			Payment:     m.Transaction.Payment,
			Status:      st,
		}
	}
	return setHistoryOrder
}

func SetDetailResponseOrderHistory(order model.Order) SetDetailPaymentHistory {
	st := order.Transaction.Status
	if st == "PAID" {
		st = "Payment Approved"
	} else {
		st = "Payment Not Approved"
	}

	subtotal := order.Transaction.Amount / order.Transaction.Service // Menghitung subtotal dengan membagi jumlah amount dengan 10%

	res := SetDetailPaymentHistory{
		ID:          order.ID,
		OrderID:     order.OrderCode,
		Name:        order.Name,
		TableNumber: order.NumberTable,
		Item:        []SetItemDetailPaymentHistory{},
		SubTotal:    subtotal,
		Service:     order.Transaction.Service,
		GrandTotal:  order.Transaction.Amount,
		Payment:     order.Transaction.Payment,
		Status:      st,
	}

	for _, item := range order.Items {
		setItem := SetItemDetailPaymentHistory{
			Quantity: item.Quantity,
			Name:     item.ProductName,
			Price:    item.Subtotal,
		}
		res.Item = append(res.Item, setItem)
	}

	return res
}
