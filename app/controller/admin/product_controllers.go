package admin

import (
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetProductsController(c echo.Context) error {
	var products []model.Products

	if err := config.Db.Find(&products).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all Product",
		"product": products,
	})
}
func CreateProductController(c echo.Context) error {
	product := model.Products{}
	c.Bind(&product)

	if err := config.Db.Save(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create new product",
		"product": product,
	})
}

func GetProductController(c echo.Context) error {
	var product model.Products

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := config.Db.First(&product, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get product by id",
		"product": product,
	})
}

func UpdateProductController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var product model.Products
	if err := config.Db.First(&product, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	c.Bind(&product)

	if err := config.Db.Save(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success update product by id",
		"product": product,
	})
}
func DeleteProductController(c echo.Context) error {
	id := c.Param("id")

	product := model.Products{}
	err := config.Db.Where("id = ?", id).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve product")
	}

	if err := config.Db.Delete(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete product")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product deleted successfully",
	})

}
