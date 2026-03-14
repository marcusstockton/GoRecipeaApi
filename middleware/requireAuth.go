package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"recipea.com/m/database"
	domain "recipea.com/m/shared"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("In Middleware")

	// Get the token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"error": "Authorization header is required",
		})
		c.Abort()
		return
	}
	// Remove the "Bearer " prefix if it exists
	if tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	var hmacSampleSecret = []byte("your-secret-key") // TODO: move this to an environment variable. Remove duplicates in routes.go
	// decode the token and validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["exp"], claims["sub"])
		// check the expiration time
		if time.Now().Unix() > int64(claims["exp"].(float64)) {
			c.JSON(401, gin.H{
				"error": "Token has expired",
			})
			c.Abort()
			return
		}

		// attach the user ID to the request context
		var user domain.Chef
		result := database.DB.Table("chefs").Where("id = ?", claims["sub"]).First(&user)

		if result.Error != nil || user.ID == 0 {
			c.JSON(401, gin.H{
				"error": "Invalid token: user not found. Error: " + result.Error.Error(),
			})
			c.Abort()
			return

		}
		c.Set("user", user)

		// continue to the next handler
		c.Next()
	} else {
		fmt.Println(err)
	}
}
