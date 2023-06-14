package routes

import (
	middleware2 "github.com/labstack/echo/v4/middleware"
	"point-of-sale/app/controller"
	"point-of-sale/app/controller/admin"
	"point-of-sale/app/middleware"

	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo) {
	e.Use(middleware2.CORS())
	api := e.Group("api/v1")
	api.Static("/images", "./images")
	//Role cashier
	RouteCashier := api.Group("/cashier")
	RouteCashier.POST("/login", controller.LoginCashier)
	RouteCashier.Use(middleware.JWTMiddleware)
	{
		RouteCashier.GET("/order", controller.SearchItems)
		RouteCashier.GET("/order/search", controller.SearchItemsByName)
		RouteCashier.GET("/order/item/:id", controller.GetItemsByID)
		RouteCashier.POST("/checkout", controller.RequestPayment)

		RouteCashier.GET("/order/category", admin.IndexCategory)
		RouteCashier.POST("/membership", admin.AddMembership)
	}

	RouteAdmin := api.Group("/admin")
	RouteAdmin.POST("/login", controller.LoginAdmin)
	RouteAdmin.Use(middleware.JWTMiddleware, middleware.AdminMiddleware)
	{
		RouteAdmin.GET("/cashier", admin.GetCashier)
		RouteAdmin.POST("/cashier", admin.AddCashier)
		RouteAdmin.PUT("/cashier/:id", admin.EditCashier)
		RouteAdmin.DELETE("/cashier/:id", admin.DeleteCashier)
		RouteAdmin.GET("/cashier/search", admin.GetCashierByUserCode)

		RouteAdmin.GET("/membership", admin.GetMembership)
		RouteAdmin.POST("/membership", admin.AddMembership)
		RouteAdmin.POST("/membership/point", admin.AddPoint)
		RouteAdmin.PUT("/membership/:id", admin.EditMembership)
		RouteAdmin.DELETE("/membership/:id", admin.DeleteMembership)
		RouteAdmin.GET("/membership/search", admin.SearchMembership)

		RouteAdmin.GET("/orders", admin.IndexOrder)
		RouteAdmin.GET("/orders/:id", admin.DetailOrder)

		RouteAdmin.GET("/product", admin.IndexProducts)
		RouteAdmin.GET("/product/:id", admin.DetailProducts)
		RouteAdmin.POST("/product", admin.CreateProducts)
		RouteAdmin.DELETE("/product/:id", admin.DeleteProducts)
		RouteAdmin.PUT("/product/:id", admin.UpdateProducts)

		RouteAdmin.GET("/category", admin.IndexCategory)
		RouteAdmin.POST("/category", admin.CreateCategory)
		RouteAdmin.DELETE("/category/:id", admin.DeleteCategory)
	}
}
