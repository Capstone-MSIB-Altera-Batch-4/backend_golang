package res

type FormatApi struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type FormatApiPage struct {
	Meta       Meta        `json:"meta"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func Response(code int, status, message string, data interface{}) FormatApi {
	Meta := Meta{
		Code:    code,
		Status:  status,
		Message: message,
	}

	return FormatApi{
		Meta: Meta,
		Data: data,
	}
}

func Responsedata(code int, status, message string, data interface{}, pages Pagination) FormatApiPage {
	meta := Meta{
		Code:    code,
		Status:  status,
		Message: message,
	}

	pagination := Pagination{
		Page:       pages.Page,
		Limit:      pages.Limit,
		TotalItems: pages.TotalItems,
		TotalPages: (pages.TotalItems + pages.Limit - 1) / pages.Limit,
	}

	return FormatApiPage{
		Meta:       meta,
		Data:       data,
		Pagination: pagination,
	}
}
