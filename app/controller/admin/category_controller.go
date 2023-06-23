package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/dto"
	generator "point-of-sale/utils/gen"
	"point-of-sale/utils/res"
	"strconv"
)

func IndexCategory(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	var category []model.Category
	query := config.Db
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&category).Error; err != nil {
		format := res.Response(http.StatusInternalServerError, "error", "error retrieving data", err.Error())
		return c.JSON(http.StatusInternalServerError, format)
	}

	totalItems := len(category)
	pages := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
	}

	categories := res.TransformCategory(category)
	format := res.Responsedata(http.StatusOK, "success", "successfully retrieved data", categories, pages)
	return c.JSON(http.StatusOK, format)
}

func CreateCategory(c echo.Context) error {
	request := dto.CreateCategoryRequest{}
	if err := c.Bind(&request); err != nil {
		format := res.Response(http.StatusInternalServerError, "error", "error request body", err.Error())
		return c.JSON(http.StatusInternalServerError, format)
	}

	if err := generator.ValidateData(&request); err != nil {
		response := res.Response(http.StatusBadRequest, "error", "failed input data", err)
		return c.JSON(http.StatusBadRequest, response)
	}

	category := model.Category{
		Name: request.Name,
	}

	if err := config.Db.Create(&category).Error; err != nil {
		format := res.Response(http.StatusInternalServerError, "error", "error create data", err.Error())
		return c.JSON(http.StatusInternalServerError, format)
	}
	format := res.Response(http.StatusCreated, "success", "Category created successfully", category)
	return c.JSON(http.StatusCreated, format)
}

func DeleteCategory(c echo.Context) error {
	categoryID := c.Param("id")

	if categoryID == "" {
		format := res.Response(http.StatusBadRequest, "error", "Invalid product ID", nil)
		return c.JSON(http.StatusBadRequest, format)
	}

	id, err := strconv.Atoi(categoryID)
	if err != nil {
		format := res.Response(http.StatusBadRequest, "error", "Invalid product ID", nil)
		return c.JSON(http.StatusBadRequest, format)
	}

	if err := config.Db.Delete(&model.Category{}, id).Error; err != nil {
		format := res.Response(http.StatusInternalServerError, "error", "error delete data", err.Error())
		return c.JSON(http.StatusInternalServerError, format)
	}

	format := res.Response(http.StatusOK, "success", "CategoryID deleted successfully", nil)
	return c.JSON(http.StatusOK, format)
}
