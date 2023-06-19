package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/res"
	"strconv"
)

func OrderHistory(c echo.Context) error {
	current := c.Get("user").(model.User)
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var orders []model.Order
	var count int64
	result := config.Db.Model(&model.Order{}).Preload("Items").Preload("Transaction").Joins("Transaction").
		Where("user_id = ?", current.ID).
		//kunci Count - Offset - limit
		Count(&count).Offset(offset).Limit(limit).
		Order("orders.created_at DESC").
		Find(&orders)

	pagination := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(count),
	}

	if result.Error != nil {
		response := res.Response(http.StatusInternalServerError, "error", "failed to retrieve order history", result.Error.Error())
		return c.JSON(http.StatusInternalServerError, response)
	}
	format := res.SetResponseOrderHistory(orders)
	response := res.Responsedata(http.StatusOK, "success", "successfully retrieved order history", format, pagination)
	return c.JSON(http.StatusOK, response)
}

func DetailOrderHistory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		format := res.Response(http.StatusBadRequest, "error", "Invalid Order History ID", nil)
		return c.JSON(http.StatusBadRequest, format)
	}

	current := c.Get("user").(model.User)

	var orders model.Order
	result := config.Db.Preload("Items").Preload("Transaction").Where("id = ?", id).Find(&orders)
	if result.Error != nil {
		response := res.Response(http.StatusInternalServerError, "error", "failed to retrieve order history", result.Error.Error())
		return c.JSON(http.StatusInternalServerError, response)
	}

	if orders.Transaction.UserID != current.ID {
		response := res.Response(http.StatusInternalServerError, "error", "failed to retrieve order history", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}
	format := res.SetDetailResponseOrderHistory(orders)
	response := res.Response(http.StatusOK, "success", "successfully retrieved order history", format)
	return c.JSON(http.StatusOK, response)
}
