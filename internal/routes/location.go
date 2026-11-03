package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

func LocationRoutes(g *echo.Group, h *handlers.LocationHandler) {
	g.POST("/locations", h.Create)
	g.GET("/locations", h.GetAll)
	g.GET("/locations/:id", h.GetByID)
	g.PUT("/locations/:id", h.Update)
	g.DELETE("/locations/:id", h.Delete)

	// Routes Many-to-Many
	g.POST("/locations/:id/entities/:entity_id", h.AddEntity)
	g.DELETE("/locations/:id/entities/:entity_id", h.RemoveEntity)
	g.POST("/locations/:id/items/:item_id", h.AddItem)
	g.DELETE("/locations/:id/items/:item_id", h.RemoveItem)
}
