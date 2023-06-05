package dto

type AddCashierRequest struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

type EditCashierRequest struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}


