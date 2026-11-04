package routes

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
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

	// 3. Creating base group (public) and protected group (admin)
	api := e.Group("/api/v1")
	admin := api.Group("")
	admin.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	// 4. Register routes
	AuthRoutes(api)

	EntityRoutes(api, admin, entityHandler)
	LocationRoutes(api, admin, locationHandler)
	ItemRoutes(api, admin, itemHandler)
}
