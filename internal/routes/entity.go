package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

// EntityRoutes all routes about Entity
func EntityRoutes(public *echo.Group, admin *echo.Group, h *handlers.EntityHandler) {
	// Public routes (Read-only)
	public.GET("/entities", h.GetAll)
	public.GET("/entities/:id", h.GetByID)

	// Admin routes (Protected by JWT)
	admin.POST("/entities", h.Create)
	admin.PUT("/entities/:id", h.Update)
	admin.DELETE("/entities/:id", h.Delete)

	// Routes Many-to-Many (Admin only)
	admin.POST("/entities/:id/drops/:item_id", h.AddDrop)
	admin.DELETE("/entities/:id/drops/:item_id", h.RemoveDrop)
}
