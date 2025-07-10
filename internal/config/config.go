package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
	JWTSecret   string
	Environment string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		Port:        getEnv("PORT", "4000"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBUser:      getEnv("user", ""),
		DBPassword:  getEnv("password", ""),
		DBHost:      getEnv("host", ""),
		DBPort:      getEnv("port", "5432"),
		DBName:      getEnv("dbname", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	return config, nil
}

func (c *Config) GetDSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=require",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}