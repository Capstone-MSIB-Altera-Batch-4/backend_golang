package gen

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"path/filepath"
	"strings"
)

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func ValidateData(data interface{}) []ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("image", validateImageExtension)

	if err := validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMsg := make(map[string]string)

		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field %s wajib diisi.", fieldError.StructField())
			case "email":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field %s harus berisi alamat email yang valid.", fieldError.StructField())
			case "min":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field %s harus memiliki panjang minimal %s karakter.", fieldError.StructField(), fieldError.Param())
			case "max":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field %s harus memiliki panjang maksimal %s karakter.", fieldError.StructField(), fieldError.Param())
			case "numeric":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field %s harus berisi angka.", fieldError.StructField())
			case "alpha":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field %s hanya boleh berisi karakter.", fieldError.StructField())
			case "len":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field %s harus memiliki panjang %s.", fieldError.StructField(), fieldError.Param())
			case "image":
				errorMsg[fieldError.StructField()] = fmt.Sprintf("Field %s harus memiliki ekstensi file .jpg.", fieldError.StructField())
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
	file, ok := fl.Field().Interface().(string)
	if !ok || file == "" {
		return false
	}

	ext := strings.ToLower(filepath.Ext(file))
	return ext == ".jpg" || ext == ".jpeg"
}
