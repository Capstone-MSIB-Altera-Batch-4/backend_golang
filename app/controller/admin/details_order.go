package admin

import (
	"fmt"
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

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	var orders []model.Order
	var totalItems int64
	query := config.Db.Model(&model.Order{})

	if orderCode != "" {
		query = query.Where("order_code LIKE ?", "%"+orderCode+"%")
	}

	endDate = fmt.Sprintf("%s 23:59:59", endDate)
	startDate = fmt.Sprintf("%s 00:00:00", startDate)
	if startDate != "" && endDate != "" {
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Count(&totalItems)

	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Preload("Transaction").Find(&orders).Error; err != nil {
		response := res.Response(500, "error", "Internal Server Error", err.Error())
		return c.JSON(500, response)
	}

	var transformedOrders []res.SetOrderResponse
	for _, order := range orders {
		transformedOrder := res.TransformResponseDataOrder(order)
		transformedOrders = append(transformedOrders, transformedOrder)
	}

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
