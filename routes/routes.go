package routes

import (
	"gin/config"
	"gin/controllers"
	"gin/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(cfg *config.Config) *gin.Engine {
    r := gin.Default()

    // Middleware
    r.Use(middleware.CORSMiddleware())

    // Initialize controllers
    authController := controllers.NewAuthController(cfg)
    userController := controllers.NewUserController()
    companyController := controllers.NewCompanyController()
    traineeController := controllers.NewTraineeController()
    timeRecordController := controllers.NewTimeRecordController()
    reportsController := controllers.NewReportsController()

    // Public routes
    auth := r.Group("/api/auth")
    {
        auth.POST("/register", authController.Register)
        auth.POST("/login", authController.Login)
        auth.POST("/logout", authController.Logout)
    }

    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware(cfg))
    {
        // User routes
        api.GET("/profile", userController.GetProfile)
        api.GET("/users", userController.GetUsers)

        // Company routes
        companies := api.Group("/companies")
        {
            companies.POST("/", companyController.CreateCompany)
            companies.GET("/", companyController.GetCompanies)
            companies.GET("/:id", companyController.GetCompany)
            companies.PUT("/:id", companyController.UpdateCompany)
            companies.DELETE("/:id", companyController.DeleteCompany)
        }

        // Trainee routes
        trainees := api.Group("/trainees")
        {
            trainees.POST("/", traineeController.CreateTrainee)
            trainees.GET("/", traineeController.GetTrainees)
            trainees.GET("/:id", traineeController.GetTrainee)
            trainees.PUT("/:id", traineeController.UpdateTrainee)
            trainees.DELETE("/:id", traineeController.DeleteTrainee)
        }

        // Time record routes
        timeRecords := api.Group("/time-records")
        {
            timeRecords.POST("/", timeRecordController.CreateTimeRecord)
            timeRecords.GET("/", timeRecordController.GetTimeRecords)
            timeRecords.GET("/:id", timeRecordController.GetTimeRecord)
            timeRecords.PUT("/:id", timeRecordController.UpdateTimeRecord)
            timeRecords.DELETE("/:id", timeRecordController.DeleteTimeRecord)
        }

        // Reports routes
        reports := api.Group("/reports")
        {
            reports.GET("/weekly", reportsController.GetWeeklySummary)
            reports.GET("/monthly", reportsController.GetMonthlyReport)
            reports.GET("/dtr", reportsController.GetDTR)
            reports.GET("/ojt-progress/:trainee_id", reportsController.GetOJTProgress)
            reports.GET("/ojt-progress", reportsController.GetAllOJTProgress)
        }
    }

    return r
}