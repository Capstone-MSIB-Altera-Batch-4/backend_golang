package routes

import (
	"github.com/labstack/echo/v4"
	"point-of-sale/app/controller"
	"point-of-sale/app/controller/admin"
	"point-of-sale/app/middleware"
)

func Route(e *echo.Echo) {
	api := e.Group("api/v1")

	//Role cashier
	RouteCashier := api.Group("/cashier")
	RouteCashier.POST("/login", controller.LoginCashier)
	RouteCashier.Use(middleware.JWTMiddleware)
	{
		e.GET("/order", controller.SearchItems)
		e.GET("/order/search", controller.SearchItemsByName)
		e.GET("/order/member", controller.SearchMembershipByName)
		e.POST("/checkout", controller.RequestPayment)
	}

	RouteAdmin := api.Group("/admin")
	RouteAdmin.POST("/login", controller.LoginAdmin)
	RouteAdmin.Use(middleware.JWTMiddleware, middleware.AdminMiddleware)
	{
		RouteAdmin.GET("/cashier", admin.GetCashier)
		RouteAdmin.POST("/cashier", admin.AddCashier)
		RouteAdmin.PUT("/cashier/:id", admin.EditCashier)
		RouteAdmin.DELETE("/cashier/:id", admin.DeleteCashier)

		RouteAdmin.GET("/membership", admin.GetMembership)
		RouteAdmin.POST("/membership", admin.AddMembership)
		RouteAdmin.POST("/membership/point", admin.AddPoint)
		RouteAdmin.PUT("/membership/:id", admin.EditMembership)
		RouteAdmin.DELETE("/membership/:id", admin.DeleteMembership)

		RouteAdmin.GET("/orders", admin.IndexOrder)
		RouteAdmin.GET("/orders/:id", admin.DetailOrder)

		RouteAdmin.GET("/product", admin.IndexProducts)
		RouteAdmin.GET("/product/:id", admin.DetailProducts)
		RouteAdmin.POST("/product/create", admin.CreateProducts)
		RouteAdmin.DELETE("/product/delete", admin.DeleteProducts)
		RouteAdmin.PUT("/product/update", admin.UpdateProducts)

		RouteAdmin.GET("/category", admin.IndexCategory)
		RouteAdmin.POST("/category/create", admin.CreateCategory)
		RouteAdmin.DELETE("/category/delete", admin.DeleteCategory)
	}

	e.Static("/images", "./images")
}
