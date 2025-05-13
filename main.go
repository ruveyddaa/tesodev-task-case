package main

import (
	"fmt"
	"net/http"
	"tesodev-product-api/config"
	"tesodev-product-api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	config.ConnectDB()
	routes.ProductRoutes(e)

	for _, r := range e.Routes() {
		fmt.Printf("%-6s → %s\n", r.Method, r.Path)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
