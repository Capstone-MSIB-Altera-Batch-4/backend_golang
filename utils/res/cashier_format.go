package res

import (
	"point-of-sale/app/model"
	"time"
)

type SetCashierFormat struct {
	ID        int       `json:"id"`
	UserCode  string    `json:"id_code"`
	Username  string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func TransformCashiers(request []*model.User) []SetCashierFormat {
	transformedCategory := make([]SetCashierFormat, len(request))
	for i, cashier := range request {
		transformedCategory[i] = SetCashierFormat{
			ID:        cashier.ID,
			UserCode:  cashier.UserCode,
			Username:  cashier.Username,
			Role:      cashier.Role,
			CreatedAt: cashier.CreatedAt,
			UpdatedAt: cashier.UpdatedAt,
		}
	}
	return transformedCategory
}

func TransformCashier(request model.User) SetCashierFormat {
	return SetCashierFormat{
		ID:        request.ID,
		UserCode:  request.UserCode,
		Username:  request.Username,
		Role:      request.Role,
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
	}
}
