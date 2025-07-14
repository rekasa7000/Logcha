package controllers

import (
	"gin/database"
	"gin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
    return &UserController{}
}

func (uc *UserController) GetProfile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user": models.UserResponse{
            ID:       user.ID,
            Email:    user.Email,
        },
    })
}

func (uc *UserController) GetUsers(c *gin.Context) {
    var users []models.User
    if err := database.DB.Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }

    var userResponses []models.UserResponse
    for _, user := range users {
        userResponses = append(userResponses, models.UserResponse{
            ID:       user.ID,
            Email:    user.Email,
        })
    }

    c.JSON(http.StatusOK, gin.H{"users": userResponses})
}