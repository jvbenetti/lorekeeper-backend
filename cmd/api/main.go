package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Initialize Echo instance
	e := echo.New()

	// 2. Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS is essential since your React Web/Admin panel will run on a different domain
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // You can restrict this to your Vercel/Netlify domain later
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// TODO: Initialize Database Connection (Supabase/PostgreSQL)
	// db := config.ConnectDatabase()

	// TODO: Initialize Repositories and Handlers
	// ...

	// TODO: Register Routes
	// ...

	// 3. Health Check Route
	// A simple route to verify if the API is alive
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "online",
			"world":  "Lamúria Codex API is running!",
		})
	})

	// 4. Server Port Configuration
	// Railway automatically injects the PORT environment variable.
	// If it's empty (running locally on your Mac), it defaults to 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 5. Start the server
	log.Printf("Starting server on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
