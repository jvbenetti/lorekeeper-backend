package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/repository"
)

// SetupRoutes center all routes
func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// 1. Init repos
	entityRepo := &repository.EntityRepository{DB: db}
	// Futuramente: locationRepo := &repository.LocationRepository{DB: db}
	// Futuramente: itemRepo := &repository.ItemRepository{DB: db}

	// 2. Init Handlers
	entityHandler := &handlers.EntityHandler{Repo: entityRepo}
	// Futuramente: locationHandler := &handlers.LocationHandler{Repo: locationRepo}

	// 3. Creating base group
	api := e.Group("/api/v1")

	// 4. Register routes
	EntityRoutes(api, entityHandler)
	// Futuramente: LocationRoutes(api, locationHandler)
	// Futuramente: ItemRoutes(api, itemHandler)
}
