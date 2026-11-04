package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

// AuthRoutes routes about auth
func AuthRoutes(public *echo.Group) {
	public.POST("/login", handlers.Login)
}
