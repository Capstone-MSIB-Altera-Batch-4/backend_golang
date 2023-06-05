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

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
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
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	existingCashier := model.User{}
	if err := config.Db.First(&existingCashier, intID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	updatedFields := model.User{
		Username: request.Username,
		Password: string(hash),
		Role:     request.Role,
	}

	if err := config.Db.Model(&existingCashier).Updates(updatedFields).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(200, "Success", "Cashier edited", existingCashier)

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

	response := res.Response(200, "Success", "Cashier deleted", nil)

	return c.JSON(http.StatusOK, response)
}
