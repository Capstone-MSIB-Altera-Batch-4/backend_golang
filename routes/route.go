package routes

import (
	"point-of-sale/app/controller"
	"point-of-sale/app/controller/admin"
	"point-of-sale/app/middleware"

	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo) {

	e.POST("/register", controller.Register)
	e.POST("/login", controller.Login)

	e.GET("/user", controller.Current, middleware.JWTMiddleware)
	e.GET("/admin", controller.Current, middleware.JWTMiddleware, middleware.AdminMiddleware)

	e.GET("/products", admin.GetProductsController)
	e.POST("/products", admin.CreateProductController)
	e.GET("/products/:id", admin.GetProductController)
	e.PUT("/products/:id", admin.UpdateProductController)
	e.DELETE("/products/:id", admin.DeleteProductController)

}
