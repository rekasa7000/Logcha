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
        api.GET("/profile", userController.GetProfile)
        api.GET("/users", userController.GetUsers)
    }

    return r
}