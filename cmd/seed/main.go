package main

import (
	"fmt"
	"log"
	"os"

	"honda-leasing-api/configs"
	"honda-leasing-api/internal/infrastructure/database"
	"honda-leasing-api/pkg/crypto"
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

	// 2. Init Database connection
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Fatal error initializing DB: %v", err)
	}

	// 3. Hash default password
	password := "password123"
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// 4. Seeder data (Phone numbers and emails must be unique)
	users := []struct {
		Phone    string
		Email    string
		Name     string
		RoleName string
	}{
		{
			Phone:    "+6280011112222",
			Email:    "customer@gmail.com",
			Name:     "Budi Customer",
			RoleName: "CUSTOMER",
		},
		{
			Phone:    "+6280011113333",
			Email:    "officer@honda.co.id",
			Name:     "Siti Officer",
			RoleName: "ADMIN_CABANG",
		},
		{
			Phone:    "+6280011114444",
			Email:    "delivery@honda.co.id",
			Name:     "Andi Delivery",
			RoleName: "SALES",
		},
	}

	for _, u := range users {
		// Get Role ID
		var roleID int
		err := db.Table("account.roles").Select("role_id").Where("role_name = ?", u.RoleName).Scan(&roleID).Error
		if err != nil {
			log.Printf("Failed to get role ID for %s: %v", u.RoleName, err)
			continue
		}

		// Insert User
		var userID int
		err = db.Raw(`
			INSERT INTO account.users (phone_number, email, full_name, password, is_active, created_at)
			VALUES (?, ?, ?, ?, true, CURRENT_TIMESTAMP)
			ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password
			RETURNING user_id
		`, u.Phone, u.Email, u.Name, hashedPassword).Scan(&userID).Error

		if err != nil {
			log.Printf("Failed to insert user %s: %v", u.Email, err)
			continue
		}

		// Check if user already has this role
		var roleExists bool
		err = db.Raw("SELECT EXISTS(SELECT 1 FROM account.user_roles WHERE user_id = ? AND role_id = ?)", userID, roleID).Scan(&roleExists).Error
		if err != nil {
			log.Printf("Failed to check role for user %s: %v", u.Email, err)
			continue
		}

		if !roleExists {
			// Link Role
			err = db.Exec(`
				INSERT INTO account.user_roles (user_id, role_id, assigned_by)
				VALUES (?, ?, 1)
			`, userID, roleID).Error

			if err != nil {
				log.Printf("Failed to assign role to user %s: %v", u.Email, err)
				continue
			}
		}

		// Jika role CUSTOMER, pastikan ada record di dealer.customers
		if u.RoleName == "CUSTOMER" {
			var customerExists bool
			err = db.Raw("SELECT EXISTS(SELECT 1 FROM dealer.customers WHERE user_id = ?)", userID).Scan(&customerExists).Error
			if err != nil {
				log.Printf("Failed to check customer profile for user %s: %v", u.Email, err)
				continue
			}

			if !customerExists {
				err = db.Exec(`
					INSERT INTO dealer.customers (user_id, nik, nama_lengkap, no_hp, email, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				`, userID, fmt.Sprintf("00000000000000%02d", userID), u.Name, u.Phone, u.Email).Error
				if err != nil {
					log.Printf("Failed to create customer profile for user %s: %v", u.Email, err)
					continue
				}
				fmt.Printf("   ↳ Customer profile created for %s\n", u.Email)
			}
		}

		fmt.Printf("✅ Seeded user: %s (%s) / password123\n", u.Email, u.RoleName)
	}

	fmt.Println("🎉 Seeding complete!")
}
