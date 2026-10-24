package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

// EntityRoutes all routes about Entity
func EntityRoutes(g *echo.Group, h *handlers.EntityHandler) {
	g.POST("/entities", h.Create)
	g.GET("/entities", h.GetAll)
	g.GET("/entities/:id", h.GetByID)

	// Routes Many-to-Many
	g.POST("/entities/:id/drops/:item_id", h.AddDrop)
}
