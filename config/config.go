package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	JWTSecret string
	Environment string
}

func LoadConfig() *Config{
	err := godotenv.Load()
	if err!= nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port: getEnv("PORT", "8080"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "8080"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName: getEnv("DB_NAME", "mydb"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != ""{
		return value
	}
	return defaultValue
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}