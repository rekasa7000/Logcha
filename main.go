package main

import (
	"gin/config"
	"gin/database"
	"gin/routes"
	"log"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Connect to database
    database.Connect(cfg)

    // Setup routes
    r := routes.SetupRoutes(cfg)

    // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    if err := r.Run(":" + cfg.Port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}