//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	"honda-leasing-api/configs"
	"honda-leasing-api/internal/auth"
	authHandler "honda-leasing-api/internal/auth/handler"
	authRepo "honda-leasing-api/internal/auth/postgres"
	"honda-leasing-api/internal/catalog"
	catalogHandler "honda-leasing-api/internal/catalog/handler"
	catalogRepo "honda-leasing-api/internal/catalog/postgres"
	"honda-leasing-api/internal/delivery"
	deliveryHandler "honda-leasing-api/internal/delivery/handler"
	deliveryRepo "honda-leasing-api/internal/delivery/postgres"
	"honda-leasing-api/internal/infrastructure/http"
	"honda-leasing-api/internal/infrastructure/http/swagger"
	"honda-leasing-api/internal/leasing"
	leasingHandler "honda-leasing-api/internal/leasing/handler"
	leasingRepo "honda-leasing-api/internal/leasing/postgres"
	"honda-leasing-api/internal/middleware"
	"honda-leasing-api/internal/officer"
	officerHandler "honda-leasing-api/internal/officer/handler"
	officerRepo "honda-leasing-api/internal/officer/postgres"
)

// ConfigProvider provides specific config parts extracted from the global config
var ConfigProviderSet = wire.NewSet(
	ProvideJwtConfig,
)

func ProvideJwtConfig(cfg *configs.Config) configs.JwtConfig {
	return cfg.Jwt
}

var AuthProviderSet = wire.NewSet(
	authRepo.NewAuthRepository,
	auth.NewService,
	authHandler.NewAuthHandler,
)

var CatalogProviderSet = wire.NewSet(
	catalogRepo.NewCatalogRepository,
	catalog.NewService,
	catalogHandler.NewCatalogHandler,
)

var LeasingProviderSet = wire.NewSet(
	leasingRepo.NewLeasingRepository,
	leasing.NewService,
	leasingHandler.NewLeasingHandler,
)

var OfficerProviderSet = wire.NewSet(
	officerRepo.NewOfficerRepository,
	officer.NewService,
	officerHandler.NewOfficerHandler,
)

var DeliveryProviderSet = wire.NewSet(
	deliveryRepo.NewDeliveryRepository,
	delivery.NewService,
	deliveryHandler.NewDeliveryHandler,
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.Auth,
)

// Server setup that registers all routes and returns the server instance
func ProvideServer(
	cfg *configs.Config,
	authHTTPHandler *authHandler.AuthHandler,
	catHandler *catalogHandler.CatalogHandler,
	leasHandler *leasingHandler.LeasingHandler,
	offcHandler *officerHandler.OfficerHandler,
	delivHandler *deliveryHandler.DeliveryHandler,
	authMiddleware gin.HandlerFunc,
) *http.Server {
	srv := http.NewServer(cfg.App.Port, cfg.App.Env)

	// Register global middleware
	srv.Router.Use(middleware.RequestLogger())
	srv.Router.Use(middleware.GlobalRecovery())

	// Register Swagger UI
	swagger.RegisterRoutes(srv.Router)

	// Register routes
	authHTTPHandler.RegisterRoutes(srv.Router, authMiddleware)
	catHandler.RegisterRoutes(srv.Router, authMiddleware)
	leasHandler.RegisterRoutes(srv.Router, authMiddleware, middleware.RoleBasedAccessControl)
	offcHandler.RegisterRoutes(srv.Router, authMiddleware, middleware.RoleBasedAccessControl)
	delivHandler.RegisterRoutes(srv.Router, authMiddleware, middleware.RoleBasedAccessControl)

	return srv
}

// InitializeServer establishes the dependency graph and returns the HTTP server
func InitializeServer(db *gorm.DB, cfg *configs.Config) (*http.Server, error) {
	wire.Build(
		ConfigProviderSet,
		AuthProviderSet,
		CatalogProviderSet,
		LeasingProviderSet,
		OfficerProviderSet,
		DeliveryProviderSet,
		MiddlewareProviderSet,
		ProvideServer,
	)
	return nil, nil
}
