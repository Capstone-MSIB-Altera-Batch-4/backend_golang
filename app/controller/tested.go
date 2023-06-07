package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type TestValidation struct {
	Name  string `json:"name" validate:"required"`
	Nik   string `json:"nik" validate:"required,min=3,max=5"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,numeric"`
	Image string `json:"image" validate:"required,customImageExtension"`
}

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func Validation(c echo.Context) error {
	person := new(TestValidation)
	if err := c.Bind(person); err != nil {
		return err
	}

	err := validateData(person)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Data berhasil dibuat")
}

func validateData(data interface{}) []ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("customImageExtension", validateImageExtension)

	if err := validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMsg := make(map[string]string)

		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				errorMsg[fieldError.StructField()] = "Field ini wajib diisi."
			case "email":
				errorMsg[fieldError.StructField()] = "Field ini harus berisi alamat email yang valid."
			case "min":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field ini harus memiliki panjang minimal %s karakter.", fieldError.Param())
			case "max":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field ini harus memiliki panjang maksimal %s karakter.", fieldError.Param())
			case "customImageExtension":
				errorMsg[fieldError.StructField()] = "Field ini harus memiliki ekstensi file .jpg"
			case "numeric":
				errorMsg[fieldError.StructField()] = "Field ini harus berisi angka."
			}
		}

		response := ErrorResponse{
			Errors: errorMsg,
		}

		return []ErrorResponse{response}
	}

	return nil
}

func validateImageExtension(fl validator.FieldLevel) bool {
	fileName := fl.Field().String()
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".jpg"
}
