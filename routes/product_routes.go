package routes

import (
	"tesodev-product-api/handlers"

	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Echo, h *handlers.ProductHandler) {
	e.GET("/products", h.GetAllProducts)
	e.POST("/products", h.CreateProduct)
	e.GET("/products/search", h.SearchProducts)
	e.GET("/products/:id", h.GetProductByID)
	e.PUT("/products/:id", h.UpdateProduct)
	e.PATCH("/products/:id", h.PatchProduct)
	e.DELETE("/products/:id", h.DeleteProduct)
}
