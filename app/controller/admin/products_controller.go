package admin

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/dto"
	"point-of-sale/utils/gen"
	"point-of-sale/utils/res"
	"strconv"
)

func IndexProducts(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	pageStr := c.QueryParam("page")
	searchCode := c.QueryParam("code")
	searchName := c.QueryParam("name")
	categoryStr := c.QueryParam("category")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	query := config.Db.Preload("Category")

	if searchCode != "" {
		query = query.Where("product_id LIKE ?", "%"+searchCode+"%")
	}

	if searchName != "" {
		query = query.Where("name LIKE ?", "%"+searchName+"%")
	}

	if categoryStr != "" {
		category := model.Category{}
		if err := config.Db.Where("name = ?", categoryStr).First(&category).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		query = query.Where("category_id = ?", category.ID)
	}

	var count int64
	query.Model(&model.Product{}).Count(&count)
	offset := (page - 1) * limit

	var products []model.Product
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	transformedProducts := res.TransformAdminProducts(products)

	pagination := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(count),
	}

	response := res.Responsedata(http.StatusOK, "success", "successfully retrieved data", transformedProducts, pagination)
	return c.JSON(http.StatusOK, response)
}

func CreateProducts(c echo.Context) error {
	request := dto.CreateProductRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// #upload file & validation image
	file, err := c.FormFile("products_image")
	if err != nil {
		if err.Error() == "http: no such file" {
			// Mengabaikan jika tidak ada file yang masuk
			request.ProductsImage = ""
		} else {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	} else {
		request.ProductsImage = file.Filename
	}

	// Lakukan validasi data
	if validationErrors := gen.ValidateData(request); validationErrors != nil {
		if request.ProductsImage != "" {
			err := os.Remove(request.ProductsImage)
			if err != nil {
				fmt.Printf("Gagal menghapus file image: %v\n", err)
			}
		}
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	fileReader, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer fileReader.Close()

	filename := uuid.NewString() + filepath.Ext(file.Filename)
	savePath := filepath.Join("images", "products", filename)

	dst, err := os.Create(savePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	_, err = io.Copy(dst, fileReader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	category := model.Category{}
	if err := config.Db.Where("id = ?", request.CategoryID).First(&category).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	product := model.Product{
		Name:        request.Name,
		Image:       savePath,
		ProductID:   request.ProductID,
		CategoryID:  category.ID,
		Quantity:    request.Quantity,
		Unit:        request.Unit,
		Price:       request.Price,
		Description: request.Description,
	}
	product.Category = category
	if err := config.Db.Create(&product).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	transformedProduct := res.TransformAdminProduct(product)
	format := res.Response(http.StatusCreated, "success", "Added product successfully", transformedProduct)
	return c.JSON(http.StatusCreated, format)
}

func DeleteProducts(c echo.Context) error {
	productID := c.Param("id")

	if productID == "" {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	id, err := strconv.Atoi(productID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	product := model.Product{}
	if err := config.Db.Where("id = ?", id).First(&product).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Product not found")
	}

	_ = os.Remove(product.Image)

	if err := config.Db.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	format := res.Response(http.StatusOK, "success", "Product deleted successfully", nil)
	return c.JSON(http.StatusOK, format)
}

func UpdateProducts(c echo.Context) error {
	productID := c.Param("id")

	if productID == "" {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	id, err := strconv.Atoi(productID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	product := model.Product{}
	if err := config.Db.Where("id = ?", id).First(&product).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Product not found")
	}

	request := dto.UpdateProductRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	category := model.Category{}
	if err := config.Db.Where("id = ?", request.CategoryID).First(&category).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Category not found")
	}
	product.Category = category
	updatedProduct := model.Product{
		ID:          product.ID,
		Name:        request.Name,
		ProductID:   request.ProductID,
		CategoryID:  request.CategoryID,
		Category:    category,
		Quantity:    request.Quantity,
		Unit:        request.Unit,
		Price:       request.Price,
		Description: request.Description,
		CreatedAt:   product.CreatedAt,
	}

	var isImageUploaded bool
	file, err := c.FormFile("products_image")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	//Validation image
	if file == nil {
		request.ProductsImage = product.Image
	} else {
		request.ProductsImage = file.Filename
	}
	if validationErrors := gen.ValidateData(request); validationErrors != nil {
		if request.ProductsImage != "" {
			if err != nil {
				fmt.Printf("Gagal menghapus file image: %v\n", err)
			}
		}
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	if file != nil {
		_ = os.Remove(product.Image)

		fileReader, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer fileReader.Close()

		filename := uuid.NewString() + filepath.Ext(file.Filename)
		savePath := filepath.Join("images", "products", filename)

		err = os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// Simpan file ke path yang ditentukan
		dst, err := os.Create(savePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer dst.Close()

		_, err = io.Copy(dst, fileReader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		updatedProduct.Image = savePath
		isImageUploaded = true
	}

	if !isImageUploaded {
		updatedProduct.Image = product.Image
	}

	if err := config.Db.Save(&updatedProduct).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	transformedProduct := res.TransformAdminProduct(updatedProduct)
	format := res.Response(http.StatusOK, "success", "Product updated successfully", transformedProduct)
	return c.JSON(http.StatusOK, format)
}

func DetailProducts(c echo.Context) error {
	idProduct, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	product := model.Product{}
	if err := config.Db.Preload("Category").Where("id = ?", idProduct).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Product not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	transformedProduct := res.TransformAdminProduct(product)
	format := res.Response(http.StatusOK, "success", "successfully retrieved data", transformedProduct)
	return c.JSON(http.StatusOK, format)
}
