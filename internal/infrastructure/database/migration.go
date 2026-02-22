package database

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"gorm.io/gorm"
)

// RunMigrations executes all .sql files found in the given directory using the GORM db connection.
func RunMigrations(db *gorm.DB, migrationsDir string) error {
	log.Printf("Starting database migrations from: %s", migrationsDir)

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	// Ensure files are executed in the correct order based on their names (e.g. 01-..., 02-...)
	sort.Strings(sqlFiles)

	for _, fileName := range sqlFiles {
		filePath := filepath.Join(migrationsDir, fileName)
		log.Printf("Executing migration: %s", fileName)

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading migration file %s: %v", fileName, err)
			return err
		}

		if err := db.Exec(string(content)).Error; err != nil {
			log.Printf("Error executing migration %s: %v", fileName, err)
			return err
		}
		log.Printf("Migration %s executed successfully.", fileName)
	}

	log.Println("All migrations applied successfully.")
	return nil
}
