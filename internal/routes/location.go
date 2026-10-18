package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

func LocationRoutes(g *echo.Group, h *handlers.LocationHandler) {
	g.POST("/locations", h.Create)
	g.GET("/locations/:id", h.GetByID)

	// Routes Many-to-Many
	g.POST("/locations/:id/entities/:entity_id", h.AddEntity)
	g.POST("/locations/:id/items/:item_id", h.AddItem)
}
