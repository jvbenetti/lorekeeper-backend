package routes

import (
	"log"
	"net/http"
	"os"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"gorm.io/gorm"

	"github.com/jvbenetti/lorekeeper-backend.git/internal/handlers"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/repository"
)

// SetupRoutes center all routes
func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("ERRO DE SEGURANÇA: JWT_SECRET não está configurado!")
	}
	if len(jwtSecret) < 32 {
		log.Println("AVISO DE SEGURANÇA: JWT_SECRET possui menos de 32 caracteres. Recomenda-se uma chave de pelo menos 256 bits.")
	}

	// 1. Init repos
	entityRepo := &repository.EntityRepository{DB: db}
	locationRepo := &repository.LocationRepository{DB: db}
	itemRepo := &repository.ItemRepository{DB: db}

	// 2. Init Handlers
	entityHandler := &handlers.EntityHandler{Repo: entityRepo}
	locationHandler := &handlers.LocationHandler{Repo: locationRepo}
	itemHandler := &handlers.ItemHandler{Repo: itemRepo}

	// 3. Creating base group (public) with general rate limiting and protected group (admin)
	api := e.Group("/api/v1")
	api.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(10), // 10 req/sec
				Burst:     30,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Limite de requisições excedido. Tente novamente em instantes.",
			})
		},
	}))

	admin := api.Group("")
	admin.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(jwtSecret),
		SigningMethod: "HS256",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Token inválido, ausente ou expirado",
			})
		},
	}))

	// 4. Register routes
	AuthRoutes(api)

	EntityRoutes(api, admin, entityHandler)
	LocationRoutes(api, admin, locationHandler)
	ItemRoutes(api, admin, itemHandler)
}
