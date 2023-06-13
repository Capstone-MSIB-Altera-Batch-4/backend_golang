package admin

import (
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/res"
	"strconv"

	"github.com/labstack/echo/v4"
)

func IndexOrder(c echo.Context) error {
	orderCode := c.QueryParam("order_id")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	limitStr := c.QueryParam("limit")
	pageStr := c.QueryParam("page")

	// default limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Nilai default jika tidak ada, tidak valid, atau terjadi kesalahan konversi
	}

	// default page
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 //Nilai default jika tidak ada, tidak valid atau terjadi kesalahan
	}

	var orders []model.Order
	var totalItems int64
	query := config.Db.Model(&model.Order{})

	//kondisi pencarian berdasarkan order_id
	if orderCode != "" {
		query = query.Where("order_id LIKE ?", "%"+orderCode+"%")
	}

	//	kondisi pencarian range tanggal
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	// Menghitung total items
	query.Count(&totalItems)

	//menghitung offset berdasarkan halaman saat ini
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Preload("Transaction").Find(&orders).Error; err != nil {
		response := res.Response(500, "error", "Internal Server Error", err.Error())
		return c.JSON(500, response)
	}

	// Transform orders to the desired response format
	var transformedOrders []res.SetOrderResponse
	for _, order := range orders {
		transformedOrder := res.TransformResponseDataOrder(order)
		transformedOrders = append(transformedOrders, transformedOrder)
	}

	// Calculate total pages
	pagination := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(totalItems),
	}
	response := res.Responsedata(http.StatusOK, "success", "successfully get data order", transformedOrders, pagination)

	return c.JSON(200, response)
}

func DetailOrder(c echo.Context) error {
	ID := c.Param("id")

	var order model.Order
	if err := config.Db.Preload("Items").Preload("Transaction").Where("id = ?", ID).First(&order).Error; err != nil {
		response := res.Response(404, "error", "Order not found", err.Error())
		return c.JSON(404, response)
	}

	transformedOrder := res.TransformResponse(order)
	response := res.FormatApi{
		Meta: res.Meta{
			Code:    200,
			Status:  "Success",
			Message: "Success Get Order Detail",
		},
		Data: transformedOrder,
	}

	return c.JSON(200, response)
}
