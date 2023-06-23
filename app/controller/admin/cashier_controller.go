package admin

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/dto"
	"point-of-sale/utils/gen"

	"point-of-sale/utils/res"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetCashier(c echo.Context) error {
	var (
		page     int
		limit    int
		offset   int
		total    int64
		cashiers []*model.User
	)

	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	if pageParam == "" {
		response := res.Response(http.StatusBadRequest, "error", "required parameter 'page'", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		response := res.Response(http.StatusBadRequest, "error", "parameter 'page' must be an integer", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if limitParam == "" {
		limit = 10 // Nilai default jika parameter 'limit' tidak diberikan
	} else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			response := res.Response(http.StatusBadRequest, "error", "parameter 'limit' must be an integer", nil)
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	offset = (page - 1) * limit

	if err := config.Db.Offset(offset).Where("role IN ('cashier', 'kepala cashier')").Limit(limit).Find(&cashiers).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	if err := config.Db.Model(&model.User{}).Where("role IN ('cashier', 'kepala cashier')").Count(&total).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	pages := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(total),
	}
	response := res.Responsedata(http.StatusOK, "success", "successfully retrieved data", cashiers, pages)

	return c.JSON(http.StatusOK, response)
}

func AddCashier(c echo.Context) error {
	request := dto.AddCashierRequest{}
	if err := c.Bind(&request); err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Periksa jika salah satu data kosong
	if request.Username == "" || request.Password == "" || request.Role == "" {
		response := res.Response(http.StatusBadRequest, "error", "Missing required data", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	userCode := fmt.Sprintf("%s-%d", gen.RandomStrGen(), gen.RandomIntGen())
	cashier := model.User{
		UserCode:  userCode,
		Username:  request.Username,
		Password:  string(hash),
		Role:      request.Role,
		CreatedAt: time.Now(),
	}

	if err := config.Db.Create(&cashier).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Menghilangkan field "Password" dari respons
	cashier.Password = ""

	format := res.TransformCashier(cashier)
	response := res.Response(201, "Success", "Cashier created", format)

	// Set status kode menjadi 201 Created
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)

	return c.JSON(http.StatusCreated, response)
}

func EditCashier(c echo.Context) error {
	request := dto.EditCashierRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	now := time.Now()
	cashier := model.User{
		ID:        intID,
		Username:  request.Username,
		Role:      request.Role,
		UpdatedAt: now,
	}

	if request.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		cashier.Password = string(hash)
	}

	if err := config.Db.Updates(&cashier).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := config.Db.First(&cashier, intID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Menghilangkan field "Password" dari respons
	cashier.Password = ""

	response := res.Response(200, "Success", "Cashier edited", cashier)

	return c.JSON(http.StatusOK, response)
}

func DeleteCashier(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Check if Transaction exists
	transaction := model.Transaction{}
	if err := config.Db.Where("user_id = ?", intID).First(&transaction).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if err := config.Db.Where("user_id = ?", transaction.UserID).Delete(&transaction).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// Check if OrderItems exist
	orderItems := []model.OrderItems{}
	if err := config.Db.Where("order_id = ?", transaction.OrderID).Delete(&orderItems).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// Check if Order exists
	order := model.Order{}
	if err := config.Db.Where("id = ?", transaction.OrderID).Delete(&order).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// Delete User
	if err := config.Db.Where("id = ?", intID).Delete(&model.User{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(http.StatusOK, "Success", "Cashier deleted", nil)
	return c.JSON(http.StatusOK, response)
}

func GetCashierByUserCode(c echo.Context) error {
	userCode := c.QueryParam("user_code")

	cashier := &model.User{}
	if err := config.Db.Where("role IN ('cashier', 'kepala cashier') AND user_code = ?", userCode).Delete(&cashier).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(http.StatusOK, "Success", "Cashier found", cashier)
	return c.JSON(http.StatusOK, response)
}
