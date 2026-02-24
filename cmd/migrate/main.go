package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"honda-leasing-api/configs"
	"honda-leasing-api/internal/infrastructure/database"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run cmd/migrate/main.go [up|down]")
	}

	direction := os.Args[1]
	if direction != "up" && direction != "down" {
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'", direction)
	}

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

	// 2. Init Database connection
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Fatal error initializing DB: %v", err)
	}

	// 3. Run Migrations
	migrationsPath := filepath.Join(".", "migrations")
	err = database.RunMigrations(db, migrationsPath, direction)
	if err != nil {
		log.Fatalf("Fatal error running migrations %s: %v", direction, err)
	}

	log.Printf("Successfully ran migrations %s", direction)
}
