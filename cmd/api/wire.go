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
	"honda-leasing-api/internal/finance"
	financeHandler "honda-leasing-api/internal/finance/handler"
	financeRepo "honda-leasing-api/internal/finance/postgres"
	"honda-leasing-api/internal/infrastructure/http"
	"honda-leasing-api/internal/infrastructure/http/swagger"
	"honda-leasing-api/internal/leasing"
	leasingHandler "honda-leasing-api/internal/leasing/handler"
	leasingRepo "honda-leasing-api/internal/leasing/postgres"
	"honda-leasing-api/internal/master"
	masterHandler "honda-leasing-api/internal/master/handler"
	masterRepo "honda-leasing-api/internal/master/postgres"
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

var MasterProviderSet = wire.NewSet(
	masterRepo.NewMasterRepository,
	master.NewService,
	masterHandler.NewMasterHandler,
)

var FinanceProviderSet = wire.NewSet(
	financeRepo.NewFinanceRepository,
	finance.NewService,
	financeHandler.NewFinanceHandler,
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
	masterHTTPHandler *masterHandler.MasterHandler,
	financeHTTPHandler *financeHandler.FinanceHandler,
	offcService officer.Service,
	financeService finance.Service,
	authMiddleware gin.HandlerFunc,
) *http.Server {
	// --- Dynamic Task Function Registry ---
	offcService.RegisterCallFunction("GeneratePaymentSchedule", financeService.GeneratePaymentSchedule)
	offcService.RegisterCallFunction("CreatePurchaseOrder", financeService.CreatePurchaseOrder)

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
	masterHTTPHandler.RegisterRoutes(srv.Router)
	financeHTTPHandler.RegisterRoutes(srv.Router)

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
		MasterProviderSet,
		FinanceProviderSet,
		MiddlewareProviderSet,
		ProvideServer,
	)
	return nil, nil
}
