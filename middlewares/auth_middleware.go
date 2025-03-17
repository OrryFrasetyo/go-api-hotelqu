package middlewares

import (
	"net/http"
	"strings"

	"github.com/OrryFrasetyo/go-api-hotelqu/utils"
	"github.com/gin-gonic/gin"
)

// JWTAuth is middleware for validation token JWT
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": true,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}
		
		// Format Authorization header should "Bearer [token]"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": true,
				"message": "Authorization header format must be Bearer [token]",
			})
			c.Abort()
			return
		}
		
		// validation token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": true,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}
		
		// Set employee ID and email to conteks for use by handler
		c.Set("employeeId", claims.Id)
		c.Set("employeeEmail", claims.Email)
		
		c.Next()
	}
}