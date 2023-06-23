package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"point-of-sale/config"
	"point-of-sale/routes"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	routes.Route(e)

	if err := e.Start(fmt.Sprintf(":%s", config.GetServer())); err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}
