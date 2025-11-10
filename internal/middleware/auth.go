package middleware

import (
	"net/http"
	"strings"
	"todo-api/internal/dto"
	"todo-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc{
	return func (c *gin.Context)  {
		authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Authorization header required"})
				c.Abort() 
				return 
			}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid authorization header format"})
			c.Abort()
			return 
		} 
		token := parts[1]
		claims, err := utils.ValidateToken(token) 
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid or expired token"})
			c.Abort()
			return 
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}

}