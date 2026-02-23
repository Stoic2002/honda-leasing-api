package main

import (
	"fmt"
	"log"
	"os"

	"honda-leasing-api/configs"
	"honda-leasing-api/internal/infrastructure/database"
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

	// 4. Initialize Dependency Injection with Wire
	srv, err := InitializeServer(db, cfg)
	if err != nil {
		log.Fatalf("Fatal error initializing server: %v", err)
	}

	// 5. Start Server
	if err := srv.Start(); err != nil {
		log.Fatalf("Fatal error starting server: %v", err)
	}
}
