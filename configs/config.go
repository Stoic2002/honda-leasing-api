package configs

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JwtConfig struct {
	Secret        string `mapstructure:"secret"`
	ExpireMinutes int    `mapstructure:"expire_minutes"`
	RefreshDays   int    `mapstructure:"refresh_days"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetEnvPrefix("HONDA")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Failed to unmarshal config: %v", err)
		return nil, err
	}

	// Validate required secrets
	if config.Jwt.Secret == "" {
		return nil, fmt.Errorf("JWT secret is required: set HONDA_JWT_SECRET environment variable")
	}

	return &config, nil
}
