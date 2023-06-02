package controller

import (
	"fmt"
	"net/http"
	"point-of-sale/app/middleware"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/dto"
	"point-of-sale/utils/res"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	request := dto.RequestAuthInput{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := model.User{
		Username: request.Username,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := config.Db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"message": "Account has been created",
	}

	return c.JSON(http.StatusOK, response)
}

func Login(c echo.Context) error {
	request := dto.RequestAuthInput{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := model.User{}
	result := config.Db.Where("username = ?", request.Username).First(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(http.StatusOK, "success", "successfully login", map[string]interface{}{
		"token": token,
	})

	fmt.Println("")
	fmt.Println(response)
	return c.JSON(http.StatusOK, response)
}

func Current(c echo.Context) error {
	user := c.Get("user").(model.User)
	response := map[string]interface{}{
		"data": user,
	}

	return c.JSON(http.StatusOK, response)
}
