package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

func LocationRoutes(g *echo.Group, h *handlers.LocationHandler) {
	g.POST("/locations", h.Create)
	g.GET("/locations/:id", h.GetByID)
}
