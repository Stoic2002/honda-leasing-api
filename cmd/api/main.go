package main

import (
	"fmt"
	"log"
	"os"

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
	"honda-leasing-api/internal/infrastructure/database"
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

func main() {
	// 1. Determine environment config
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	configPath := fmt.Sprintf("configs/app.%s.yaml", env)
	log.Printf("Loading configuration from %s", configPath)

	cfg, err := configs.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Fatal error loading config: %v", err)
	}

	log.Println("Configuration loaded successfully")

	// 2. Init Database connection
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Fatal error initializing DB: %v", err)
	}

	// 3. Run Migrations
	// migrationsPath := filepath.Join(".", "migrations")
	// err = database.RunMigrations(db, migrationsPath)
	// if err != nil {
	// 	log.Fatalf("Fatal error running migrations: %v", err)
	// }

	// 4. Init Dependencies & Handlers
	authMiddleware := middleware.Auth(cfg.Jwt)

	authRepository := authRepo.NewAuthRepository(db)
	authService := auth.NewService(authRepository, cfg.Jwt)
	authHTTPHandler := authHandler.NewAuthHandler(authService)

	catRepo := catalogRepo.NewCatalogRepository(db)
	catService := catalog.NewService(catRepo)
	catHandler := catalogHandler.NewCatalogHandler(catService)

	leasRepo := leasingRepo.NewLeasingRepository(db)
	leasService := leasing.NewService(leasRepo)
	leasHandler := leasingHandler.NewLeasingHandler(leasService)

	offcRepo := officerRepo.NewOfficerRepository(db)
	offcService := officer.NewService(offcRepo)
	offcHandler := officerHandler.NewOfficerHandler(offcService)

	delivRepo := deliveryRepo.NewDeliveryRepository(db)
	delivService := delivery.NewService(delivRepo)
	delivHandler := deliveryHandler.NewDeliveryHandler(delivService)

	// 5. Init HTTP Server
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

	// 6. Start Server
	if err := srv.Start(); err != nil {
		log.Fatalf("Fatal error starting server: %v", err)
	}
}
