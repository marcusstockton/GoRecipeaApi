package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appauth "recipea.com/m/application/auth"
	"recipea.com/m/domain/chef"
	"recipea.com/m/infrastructure/auth"
)

func RequireAuth(provider *auth.JWTProvider, repo chef.Repository) gin.HandlerFunc {
	authService := appauth.NewAuthService(provider, repo)
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		user, err := authService.Authenticate(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token or user not found"})
			c.Abort()
			return
		}
		c.Set("chef", user)
		c.Next()
	}
}
