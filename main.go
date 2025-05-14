package main

import (
	"tesodev-product-api/db"
	"tesodev-product-api/handlers"
	"tesodev-product-api/repository"
	"tesodev-product-api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Connect to database
	db.ConnectDB()

	// Initialize repository with db instance
	repo := &repository.ProductRepository{DB: db.DB}

	// Initialize handler with repository
	handler := &handlers.ProductHandler{Repo: repo}

	// Register routes with handler instance
	routes.ProductRoutes(e, handler)

	e.Logger.Fatal(e.Start(":8080"))
}
