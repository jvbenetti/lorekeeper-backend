package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

// EntityRoutes all routes about Entity
func EntityRoutes(g *echo.Group, h *handlers.EntityHandler) {
	g.POST("/entities", h.Create)
	g.GET("/entities/:id", h.GetByID)
}
