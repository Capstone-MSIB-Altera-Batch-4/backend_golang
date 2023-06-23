package controller

import (
	"net/http"
	"point-of-sale/app/middleware"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/dto"
	"point-of-sale/utils/res"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"
)

func LoginCashier(c echo.Context) error {
	var (
		request dto.LoginRequest
		user    model.User
	)

	if err := c.Bind(&request); err != nil {
		format := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, format)
	}

	if err := config.Db.Where("user_code = ?", request.Username).First(&user).Error; err != nil {
		format := res.Response(http.StatusInternalServerError, "error", "user code not found", nil)
		return c.JSON(http.StatusInternalServerError, format)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		format := res.Response(http.StatusInternalServerError, "error", "incorrect password", nil)
		return c.JSON(http.StatusInternalServerError, format)
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	format := res.TransformLoginResponse(user, token)
	resp := res.Response(http.StatusOK, "success", "successfully login", format)


	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, resp)
}

func LoginAdmin(c echo.Context) error {
	var (
		request dto.LoginRequest
		user    model.User
	)

	if err := c.Bind(&request); err != nil {
		format := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, format)
	}

	if err := config.Db.Where("username = ? AND role = ?", request.Username, "admin").First(&user).Error; err != nil {
		format := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, format)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		format := res.Response(http.StatusInternalServerError, "error", "incorrect password", map[string]string{"password": "incorrect password"})
		return c.JSON(http.StatusInternalServerError, format)
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	format := res.TransformLoginResponse(user, token)
	resp := res.Response(http.StatusOK, "success", "successfully login", format)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, resp)
}
