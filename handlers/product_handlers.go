package handlers

import (
	"net/http"
	"strconv"
	"time"

	"tesodev-product-api/models"
	"tesodev-product-api/repository"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllProducts(c echo.Context) error {
	products, err := repository.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, products)
}

func CreateProduct(c echo.Context) error {
	var product models.Product

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}

	if product.Name == "" || product.Price <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name and price required"})
	}

	product.ID = primitive.NewObjectID()
	product.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	id, err := repository.CreateProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	product.ID = id
	return c.JSON(http.StatusCreated, product)
}

func GetProductByID(c echo.Context) error {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid product ID"})
	}

	product, err := repository.GetProductByID(objID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid product ID"})
	}

	var updated models.Product
	if err := c.Bind(&updated); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}

	if updated.Name == "" || updated.Price <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name and price are required"})
	}

	err = repository.UpdateProduct(objID, updated)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product updated successfully"})
}

func PatchProduct(c echo.Context) error {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid product ID"})
	}

	var updateData map[string]interface{}
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}

	if price, ok := updateData["price"]; ok {
		if _, ok := price.(float64); !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Price must be a number"})
		}
	}

	err = repository.PatchProduct(objID, updateData)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product updated successfully"})
}

func DeleteProduct(c echo.Context) error {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid product ID"})
	}

	err = repository.DeleteProduct(objID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product deleted successfully"})
}

func SearchProducts(c echo.Context) error {
	name := c.QueryParam("name")
	minPriceStr := c.QueryParam("minPrice")
	maxPriceStr := c.QueryParam("maxPrice")
	sortOrder := c.QueryParam("sort")

	var minPrice, maxPrice float64
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "minPrice must be a number"})
		}
	}
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "maxPrice must be a number"})
		}
	}

	products, err := repository.SearchProducts(name, minPrice, maxPrice, sortOrder)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Search failed"})
	}

	return c.JSON(http.StatusOK, products)
}
