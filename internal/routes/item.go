package routes

import (
	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/labstack/echo/v4"
)

func ItemRoutes(g *echo.Group, h *handlers.ItemHandler) {
	g.POST("/items", h.Create)
	g.GET("/items", h.GetAll)
	g.GET("/items/:id", h.GetByID)
	g.PUT("/items/:id", h.Update)
	g.DELETE("/items/:id", h.Delete)
}
