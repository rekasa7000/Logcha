package middleware

import (
	"gin/config"
	"gin/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func (c *gin.Context)  {
		tokenString, err  := c.Cookie("jwt_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		
		claims, err := utils.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.SetCookie("jwt_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}