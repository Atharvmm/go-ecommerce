package routes

import (
	"go-ecommerce/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/products", handlers.GetAllProducts)
	e.GET("/products/:id", handlers.GetProductHandler)
	e.GET("/products/cached", handlers.GetCachedProducts)
	e.GET("/products/search", handlers.SearchProductHandler)

	e.POST("/products", handlers.AddProductHandler)
	e.POST("/products/upload", handlers.UploadProducts)

	e.PUT("/products/:id", handlers.UpdateProductHandler)

	e.DELETE("/products/delete/:id", handlers.DeleteProductHandler)
	e.DELETE("/products/delete", handlers.DeleteProductsByIDs)
	e.DELETE("/products/deleteAllProducts", handlers.DeleteAllProductsHandler)

}
