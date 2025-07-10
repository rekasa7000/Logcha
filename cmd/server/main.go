package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rekasa7000/Logcha/internal/config"
	"github.com/rekasa7000/Logcha/internal/handlers"
	"github.com/rekasa7000/Logcha/internal/middleware"
	"github.com/rekasa7000/Logcha/internal/models"
	"github.com/rekasa7000/Logcha/internal/repository"
	"github.com/rekasa7000/Logcha/internal/services"
	"github.com/rekasa7000/Logcha/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	var db *database.DB
	dsn := cfg.GetDSN()
	if dsn != "" {
		db, err = database.NewPostgresDB(dsn)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		// Auto-migrate database
		err = db.AutoMigrate(
			&models.User{},
			&models.Company{},
			&models.Trainee{},
			&models.TimeRecord{},
			&models.WeeklySummary{},
			&models.MonthlyReport{},
		)
		if err != nil {
			log.Fatal("Failed to migrate database:", err)
		}

		log.Println("Database connected and migrated successfully")
	} else {
		log.Println("No database URL provided, running without database")
	}

	// Initialize repositories
	var userRepo repository.UserRepository
	var traineeRepo repository.TraineeRepository
	var timeRecordRepo repository.TimeRecordRepository

	if db != nil {
		userRepo = repository.NewUserRepository(db)
		traineeRepo = repository.NewTraineeRepository(db)
		timeRecordRepo = repository.NewTimeRecordRepository(db)
	}

	// Initialize services
	var authService *services.AuthService
	var timeRecordService *services.TimeRecordService

	if userRepo != nil {
		authService = services.NewAuthService(userRepo, cfg.JWTSecret)
	}

	if timeRecordRepo != nil && traineeRepo != nil {
		timeRecordService = services.NewTimeRecordService(timeRecordRepo, traineeRepo)
	}

	// Initialize handlers
	var authHandler *handlers.AuthHandler
	var timeRecordHandler *handlers.TimeRecordHandler

	if authService != nil {
		authHandler = handlers.NewAuthHandler(authService)
	}

	if timeRecordService != nil {
		timeRecordHandler = handlers.NewTimeRecordHandler(timeRecordService)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CORSMiddleware())

	// API routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "Logcha API is running",
			"version": "1.0.0",
		})
	})

	// Setup routes only if handlers are available
	if authHandler != nil {
		// Public auth routes
		auth := api.Group("/auth")
		auth.Post("/login", authHandler.Login)
		auth.Post("/register", authHandler.Register)

		// Protected routes
		protected := api.Group("/", middleware.AuthMiddleware(cfg.JWTSecret))
		protected.Get("/me", authHandler.GetMe)

		// Time tracking routes
		if timeRecordHandler != nil {
			timeTracking := protected.Group("/time")
			timeTracking.Post("/in", timeRecordHandler.TimeIn)
			timeTracking.Post("/out", timeRecordHandler.TimeOut)
			timeTracking.Get("/records/:traineeId", timeRecordHandler.GetTimeRecords)
			timeTracking.Get("/today/:traineeId", timeRecordHandler.GetTodaysTimeRecord)
		}
	}

	// Fallback route for undefined endpoints
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route not found",
		})
	})

	// Start server
	log.Printf("Starting Logcha server on port %s", cfg.Port)
	log.Printf("Environment: %s", cfg.Environment)
	
	if dsn == "" {
		log.Println("⚠️  Warning: No database connection. Some features may not work.")
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.Port)))
}