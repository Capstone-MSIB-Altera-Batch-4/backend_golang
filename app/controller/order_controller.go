package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/res"
	"strconv"
)

func SearchItems(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	searchCategory := c.QueryParam("category")
	pageStr := c.QueryParam("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	categoryQuery := config.Db.Model(&model.Category{})

	if searchCategory != "" {
		categoryQuery = categoryQuery.Where("name = ?", searchCategory)
	}

	var totalItems int64
	var categories []model.Category
	if err := categoryQuery.Count(&totalItems).Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
	}

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit
	if endIndex > int(totalItems) {
		endIndex = int(totalItems)
	}

	var responseProducts []res.SetSearchOrderResponse
	for i := startIndex; i < endIndex; i++ {
		category := categories[i]

		productQuery := config.Db.Model(&model.Product{})
		if err := productQuery.Where("category_id = ?", category.ID).Find(&category.Products).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
		}
		setResponse := res.TransformCategoryOrder(category)
		responseProducts = append(responseProducts, setResponse)
	}

	pages := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(totalItems),
	}
	response := res.Responsedata(http.StatusOK, "success", "Data retrieved successfully", responseProducts, pages)

	return c.JSON(http.StatusOK, response)
}

func SearchItemsByName(c echo.Context) error {
	searchName := c.QueryParam("name")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	if searchName == "" {
		return SearchItems(c)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	var responseProducts []res.SetGetItemResponse
	var products []model.Product
	productQuery := config.Db.Model(&model.Product{}).Where("name LIKE ?", "%"+searchName+"%").Offset(offset).Limit(limit)
	if err := productQuery.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
	}

	var count int64
	if err := config.Db.Model(&model.Product{}).Where("name LIKE ?", "%"+searchName+"%").Count(&count).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
	}

	pages := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(count),
	}

	responseProducts = res.TransformItemOrder(products)
	response := res.Responsedata(http.StatusOK, "success", "Data retrieved successfully", responseProducts, pages)

	return c.JSON(http.StatusOK, response)
}
