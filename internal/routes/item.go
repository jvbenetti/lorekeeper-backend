package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

func ItemRoutes(public *echo.Group, admin *echo.Group, h *handlers.ItemHandler) {
	// Public routes (Read-only)
	public.GET("/items", h.GetAll)
	public.GET("/items/:id", h.GetByID)

	// Admin routes (Protected by JWT)
	admin.POST("/items", h.Create)
	admin.PUT("/items/:id", h.Update)
	admin.DELETE("/items/:id", h.Delete)
}
