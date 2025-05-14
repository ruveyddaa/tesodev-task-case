package routes

import (
	"tesodev-product-api/handlers"

	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Echo) {

	e.GET("/products", handlers.GetAllProducts)
	e.POST("/products", handlers.CreateProduct)

	e.GET("/products/search", handlers.SearchProducts)

	e.GET("/products/:id", handlers.GetProductByID)
	e.PUT("/products/:id", handlers.UpdateProduct)
	e.PATCH("/products/:id", handlers.PatchProduct)
	e.DELETE("/products/:id", handlers.DeleteProduct)
}
