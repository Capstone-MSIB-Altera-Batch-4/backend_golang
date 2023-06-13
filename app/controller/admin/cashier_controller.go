package admin

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/dto"
	"point-of-sale/utils/gen"

	"point-of-sale/utils/res"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func GetCashier(c echo.Context) error {
	var (
		page     int
		limit    = 10
		offset   int
		total    int64
		cashiers []*model.User
	)

	temp := c.QueryParam("page")

	if temp == "" {
		response := res.Response(http.StatusBadRequest, "error", "required paramter `page`", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	page, err := strconv.Atoi(temp)
	if err != nil {
		response := res.Response(http.StatusBadRequest, "error", "page must be integer", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	offset = (page - 1) * limit

	if err := config.Db.Offset(offset).Where("role IN ('cashier', 'kepala cashier')").Find(&cashiers).Error; err != nil {
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
	format := res.TransformCashiers(cashiers)
	response := res.Responsedata(http.StatusOK, "success", "successfully retrieved data", format, pages)
	return c.JSON(http.StatusOK, response)
}

func AddCashier(c echo.Context) error {
	request := dto.AddCashierRequest{}
	if err := c.Bind(&request); err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
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

	format := res.TransformCashier(cashier)
	response := res.Response(201, "Success", "Cashier created", format)
	return c.JSON(http.StatusOK, response)
}

func EditCashier(c echo.Context) error {
	request := dto.EditCashierRequest{}
	if err := c.Bind(&request); err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	now := time.Now()
	cashier := model.User{
		ID:        intID,
		Username:  request.Username,
		Password:  string(hash),
		Role:      request.Role,
		UpdatedAt: now,
	}

	if err := config.Db.Updates(&cashier).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	if err := config.Db.First(&cashier, intID).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	format := res.TransformCashier(cashier)
	response := res.Response(200, "Success", "Cashier edited", format)
	return c.JSON(http.StatusOK, response)
}

func DeleteCashier(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		response := res.Response(http.StatusBadRequest, "error", "Invalid cashier ID", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	result := config.Db.Where("id = ?", intID).Delete(&model.User{})
	if result.RowsAffected == 0 {
		response := res.Response(http.StatusNotFound, "error", "Cashier not found", nil)
		return c.JSON(http.StatusNotFound, response)
	}
	if result.Error != nil {
		response := res.Response(http.StatusInternalServerError, "error", result.Error.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := res.Response(http.StatusOK, "success", "Cashier deleted", nil)
	return c.JSON(http.StatusOK, response)
}

func GetCashierByUserCode(c echo.Context) error {
	userCode := c.QueryParam("user_code")

	cashier := &model.User{}
	if err := config.Db.Where("role IN ('cashier', 'kepala cashier') AND user_code = ?", userCode).First(&cashier).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	if cashier.ID == 0 {
		response := res.Response(http.StatusNotFound, "error", "Cashier not found", nil)
		return c.JSON(http.StatusNotFound, response)
	}

	format := res.TransformCashier(*cashier)
	response := res.Response(http.StatusOK, "success", "Cashier found", format)
	return c.JSON(http.StatusOK, response)
}
