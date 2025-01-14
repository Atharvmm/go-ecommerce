package main

import (
	"go-ecommerce/db"
	"go-ecommerce/routes"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	db.InitializeDB()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.RegisterRoutes(e)

	err := e.Start(":8080")

	if err != nil {
		log.Printf("failed to start the server: %v", err)
	}
}
