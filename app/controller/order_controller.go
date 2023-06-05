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
	// Query parameters
	limitStr := c.QueryParam("limit")          // Jumlah kategori yang ingin ditampilkan
	searchCategory := c.QueryParam("category") // Nama kategori yang ingin dicari
	pageStr := c.QueryParam("page")            // Nomor halaman

	// Konversi query parameter "limit" dan "page" menjadi tipe data int
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5 // Nilai default jika tidak ada, tidak valid, atau terjadi kesalahan konversi
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Nilai default jika tidak ada, tidak valid, atau terjadi kesalahan konversi
	}

	// Membuat query untuk mengambil kategori dan item
	categoryQuery := config.Db.Model(&model.Category{})

	if searchCategory != "" {
		categoryQuery = categoryQuery.Where("name = ?", searchCategory)
	}

	var totalItems int64
	var categories []model.Category
	if err := categoryQuery.Count(&totalItems).Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
	}

	// Menghitung indeks awal dan akhir item yang ditampilkan
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit
	if endIndex > int(totalItems) {
		endIndex = int(totalItems)
	}

	// Mengambil item untuk setiap kategori dengan batasan halaman dan limit
	var responseProducts []res.SetSearchOrderResponse
	for i := startIndex; i < endIndex; i++ {
		category := categories[i]

		// Membuat query untuk mengambil produk berdasarkan kategori
		productQuery := config.Db.Model(&model.Product{})
		if err := productQuery.Where("category_id = ?", category.ID).Find(&category.Products).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
		}
		setResponse := res.TransformCategoryOrder(category)
		responseProducts = append(responseProducts, setResponse)
	}

	// Membuat response
	pages := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(totalItems),
	}
	response := res.Responsedata(http.StatusOK, "success", "Data retrieved successfully", responseProducts, pages)

	// Mengembalikan response
	return c.JSON(http.StatusOK, response)
}

func SearchItemsByName(c echo.Context) error {
	// Query parameters
	searchName := c.QueryParam("name") // Nama item yang ingin dicari

	// Membuat variabel untuk menyimpan hasil pencarian
	var responseProducts []res.SetGetItemResponse

	// Jika parameter "name" kosong, panggil fungsi SearchItems
	if searchName == "" {
		return SearchItems(c)
	}

	// Mencari produk berdasarkan nama item
	var products []model.Product
	productQuery := config.Db.Model(&model.Product{}).Where("name LIKE ?", "%"+searchName+"%")
	if err := productQuery.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
	}

	// Transformasi produk menjadi respons yang diinginkan
	responseProducts = res.TransformItemOrder(products)

	// Membuat response
	response := res.Response(http.StatusOK, "success", "Data retrieved successfully", responseProducts)

	// Mengembalikan response
	return c.JSON(http.StatusOK, response)
}

func SearchMembershipByName(c echo.Context) error {
	searchName := c.QueryParam("name") // Nama membership yang ingin dicari

	if searchName == "" {
		return c.JSON(http.StatusOK, res.Response(http.StatusOK, "success", "No search term provided", nil))
	}

	var responseMemberships []model.Membership
	membershipQuery := config.Db.Model(&model.Membership{}).Where("name LIKE ?", "%"+searchName+"%")
	if err := membershipQuery.Find(&responseMemberships).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
	}

	formattedMemberships := res.TransformSearchOrderMember(responseMemberships)

	response := res.Response(http.StatusOK, "success", "Data retrieved successfully", formattedMemberships)

	return c.JSON(http.StatusOK, response)
}

//func SearchItems(c echo.Context) error {
//	// Query parameters
//	searchName := c.QueryParam("name")         // Nama item yang ingin dicari
//	limitStr := c.QueryParam("limit")          // Jumlah kategori yang ingin ditampilkan
//	searchCategory := c.QueryParam("category") // Nama kategori yang ingin dicari
//
//	// Konversi query parameter "limit" menjadi tipe data int
//	limit, err := strconv.Atoi(limitStr)
//	if err != nil || limit <= 0 {
//		limit = 5 // Nilai default jika tidak ada, tidak valid, atau terjadi kesalahan konversi
//	}
//
//	// Membuat query untuk mengambil kategori
//	categoryQuery := config.Db.Model(&model.CategoryID{})
//	if searchCategory != "" {
//		categoryQuery = categoryQuery.Where("name = ?", searchCategory)
//	}
//
//	var categories []model.CategoryID
//	if err := categoryQuery.Find(&categories).Error; err != nil {
//		return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
//	}
//
//	// Menghitung jumlah kategori
//	totalCategories := len(categories)
//
//	// Membuat variabel untuk menyimpan hasil pencarian
//	var responseProducts []res.SetSearchOrderResponse
//
//	// Mengambil item untuk setiap kategori
//	for _, category := range categories {
//		var products []model.Product
//		productQuery := config.Db.Model(&model.Product{}).Where("category_id = ?", category.ID)
//
//		// Menambahkan kondisi pencarian berdasarkan nama item
//		if searchName != "" {
//			productQuery = productQuery.Where("name LIKE ?", "%"+searchName+"%")
//		}
//
//		if err := productQuery.Find(&products).Error; err != nil {
//			return c.JSON(http.StatusInternalServerError, res.Response(http.StatusInternalServerError, "error", err.Error(), nil))
//		}
//
//		category.Products = products
//		setResponse := res.TransformCategoryOrder(category)
//		responseProducts = append(responseProducts, setResponse)
//	}
//
//	// Menentukan apakah halaman berikutnya ada
//	nextPage := totalCategories > limit
//
//	// Mendapatkan nomor halaman dari query parameter "page"
//	pageStr := c.QueryParam("page")
//	page, err := strconv.Atoi(pageStr)
//	if err != nil || page <= 0 {
//		page = 1 // Nilai default jika tidak ada, tidak valid, atau terjadi kesalahan konversi
//	}
//
//	// Menghitung jumlah item yang ditampilkan
//	startIndex := (page - 1) * limit
//	endIndex := startIndex + limit
//	if endIndex > totalCategories {
//		endIndex = totalCategories
//	}
//	showItems := fmt.Sprintf("%d of %d", endIndex, totalCategories)
//
//	// Membuat response
//	responseData := map[string]interface{}{
//		"nextPage": nextPage,
//		//"products":  responseProducts[startIndex:endIndex],
//		"showItems": showItems,
//	}
//	products := responseProducts[startIndex:endIndex]
//	response := res.ResponsePage(http.StatusOK, "success", "Data retrieved successfully", products, responseData)
//
//	// Mengembalikan response
//	return c.JSON(http.StatusOK, response)
//}
