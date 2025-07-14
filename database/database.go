package database

import (
	"fmt"
	"gin/config"
	"gin/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB *gorm.DB


func Connect(cfg *config.Config){
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Run migrations for all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Trainee{},
		&models.TimeRecord{},
		&models.WeeklySummary{},
		&models.MonthlyReport{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	log.Println("Database connected and migrated successfully!")
}