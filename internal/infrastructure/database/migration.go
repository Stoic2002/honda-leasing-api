package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"gorm.io/gorm"
)

// RunMigrations executes all .up.sql or .down.sql files found in the given directory using the GORM db connection.
// direction parameter must be "up" or "down".
func RunMigrations(db *gorm.DB, migrationsDir string, direction string) error {
	log.Printf("Starting database %s migrations from: %s", direction, migrationsDir)

	if direction != "up" && direction != "down" {
		return fmt.Errorf("invalid migration direction: %s. Use 'up' or 'down'", direction)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	var sqlFiles []string
	expectedExt := fmt.Sprintf(".%s.sql", direction)

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			// Check if file ends with .up.sql or .down.sql
			nameLen := len(file.Name())
			extLen := len(expectedExt)
			if nameLen >= extLen && file.Name()[nameLen-extLen:] == expectedExt {
				sqlFiles = append(sqlFiles, file.Name())
			}
		}
	}

	// Ensure files are executed in the correct order based on their names (e.g. 01-..., 02-...)
	sort.Strings(sqlFiles)

	// If migrating down, we must execute them in reverse order
	if direction == "down" {
		sort.Sort(sort.Reverse(sort.StringSlice(sqlFiles)))
	}

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
