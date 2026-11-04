package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

func LocationRoutes(public *echo.Group, admin *echo.Group, h *handlers.LocationHandler) {
	// Public routes (Read-only)
	public.GET("/locations", h.GetAll)
	public.GET("/locations/:id", h.GetByID)

	// Admin routes (Protected by JWT)
	admin.POST("/locations", h.Create)
	admin.PUT("/locations/:id", h.Update)
	admin.DELETE("/locations/:id", h.Delete)

	// Routes Many-to-Many (Admin only)
	admin.POST("/locations/:id/entities/:entity_id", h.AddEntity)
	admin.DELETE("/locations/:id/entities/:entity_id", h.RemoveEntity)
	admin.POST("/locations/:id/items/:item_id", h.AddItem)
	admin.DELETE("/locations/:id/items/:item_id", h.RemoveItem)
}
