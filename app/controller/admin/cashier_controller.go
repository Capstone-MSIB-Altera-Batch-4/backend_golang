package admin

import (
	"fmt"
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
		limit    = 10
		offset   int
		total    int64
		cashiers []*model.User
	)

	temp := c.QueryParam("page")

	if temp == "" {
		return c.JSON(http.StatusBadRequest, "required paramter `page`")
	}

	page, err := strconv.Atoi(temp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "page must be integer")
	}

	offset = (page - 1) * limit

	if err := config.Db.Offset(offset).Where("role IN ('cashier', 'kepala cashier')").Find(&cashiers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := config.Db.Model(&model.User{}).Where("role IN ('cashier', 'kepala cashier')").Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
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
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userCode := fmt.Sprintf("%s-%d", gen.RandomStrGen(), gen.RandomIntGen())

	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	cashier := model.User{
		UserCode:  userCode,
		Username:  request.Username,
		Password:  string(hashedPassword),
		Role:      request.Role,
		CreatedAt: time.Now(),
	}

	if err := config.Db.Create(&cashier).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(201, "Success", "Cashier created", cashier)

	return c.JSON(http.StatusOK, response)
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
		Password:  request.Password,
		Role:      request.Role,
		UpdatedAt: now,
	}

	if err := config.Db.Updates(&cashier).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := config.Db.First(&cashier, intID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	response := res.Response(200, "Success", "Cashier edited", cashier)

	return c.JSON(http.StatusOK, response)
}

func DeleteCashier(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := config.Db.Where("id = ?", intID).Delete(&model.User{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(200, "Success", "Cashier deleted", "")

	return c.JSON(http.StatusOK, response)
}

func GetCashierByUserCode(c echo.Context) error {
	userCode := c.QueryParam("user_code")

	cashier := &model.User{}
	if err := config.Db.Where("role IN ('cashier', 'kepala cashier') AND user_code = ?", userCode).First(&cashier).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(http.StatusOK, "Success", "Cashier found", cashier)
	return c.JSON(http.StatusOK, response)
}
