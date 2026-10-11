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
	locationRepo := &repository.LocationRepository{DB: db}
	itemRepo := &repository.ItemRepository{DB: db}

	// 2. Init Handlers
	entityHandler := &handlers.EntityHandler{Repo: entityRepo}
	locationHandler := &handlers.LocationHandler{Repo: locationRepo}
	itemHandler := &handlers.ItemHandler{Repo: itemRepo}

	// 3. Creating base group
	api := e.Group("/api/v1")

	// 4. Register routes
	EntityRoutes(api, entityHandler)
	LocationRoutes(api, locationHandler)
	ItemRoutes(api, itemHandler)
}
