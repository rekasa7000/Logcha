package controllers

import (
	"gin/config"
	"gin/database"
	"gin/models"
	"gin/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type AuthController struct {
	cfg *config.Config
}

func NewAuthController (cfg *config.Config) *AuthController{
	return &AuthController{cfg: cfg}
}

func (ac *AuthController) setCookie(c *gin.Context, token string){
	c.SetCookie(
		"jwt_token",
		token,
		int(24*time.Hour.Seconds()),
		"/",
		"",
		ac.cfg.IsProduction(),
		true,
	)
}

func (ac *AuthController) clearCookie(c *gin.Context){
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)
}

func (ac *AuthController) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if user already exists
    var existingUser models.User
    if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
        return
    }

    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Create user
    user := models.User{
        Email:        req.Email,
        PasswordHash: hashedPassword,
        Role:         req.Role,
        FirstName:    req.FirstName,
        LastName:     req.LastName,
        Phone:        req.Phone,
    }

    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, ac.cfg.JWTSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Set cookie
    ac.setCookie(c, token)

    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "user": models.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            FirstName: user.FirstName,
            LastName:  user.LastName,
            Phone:     user.Phone,
            Role:      user.Role,
            IsActive:  user.IsActive,
        },
    })
}

func (ac *AuthController) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Find user
    var user models.User
    if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Check password
    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, ac.cfg.JWTSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Set cookie
    ac.setCookie(c, token)

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "user": models.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            FirstName: user.FirstName,
            LastName:  user.LastName,
            Phone:     user.Phone,
            Role:      user.Role,
            IsActive:  user.IsActive,
        },
    })
}

func (ac *AuthController) Logout(c *gin.Context) {
    // Clear the JWT cookie
    ac.clearCookie(c)
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Logout successful",
    })
}