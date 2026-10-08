package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/config"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/routes"
)

func main() {
	e := echo.New()

	// Global Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// 1. Init DB
	db := config.ConnectDatabase()

	// 2. Call the setup routes
	routes.SetupRoutes(e, db)

	// Basic health check route /api/v1
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "online",
			"world":  "lorekeeper api",
		})
	})

	// Get Railway port or use 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
