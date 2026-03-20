package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"recipea.com/m/domain/chef"
	"recipea.com/m/infrastructure/auth"
)

func RequireAuth(provider *auth.JWTProvider, repo chef.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := provider.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		user, err := repo.FindByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}
		c.Set("chef", user)
		c.Next()
	}
}
