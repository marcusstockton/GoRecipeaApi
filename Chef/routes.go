package Chef

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"recipea.com/m/database"
	"recipea.com/m/domain"
	"recipea.com/m/middleware"
)

type Handler struct {
	DB *gorm.DB
}

var hmacSampleSecret = []byte("your-secret-key") // TODO: move this to an environment variable. Remove duplicate in requireAuth.go

func Routes(route *gin.Engine) {
	handler := &Handler{DB: database.DB}
	chef := route.Group("/chef")
	{
		chef.GET("/", handler.getChefs)
		chef.GET("/:id", handler.getChef)
		chef.POST("/", handler.createChef)
		chef.PUT("/:id", handler.updateChef)
		chef.DELETE("/:id", handler.deleteChef)
		chef.POST("/login", handler.Login)
		chef.GET("/validate", middleware.RequireAuth, handler.Validate)
	}
}

func (h *Handler) getChefs(c *gin.Context) {

	var chefs []domain.Chef
	h.DB.Find(&chefs)
	c.IndentedJSON(http.StatusOK, chefs)
}

func (h *Handler) createChef(c *gin.Context) {
	// Parse the request body into a Chef struct
	var body domain.Chef
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// create the chef in the database
	chef := domain.Chef{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  string(hashedPassword),
	}

	result := h.DB.Create(&chef)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "Failed to create chef",
		})
		return
	}

	// respond with the created chef
	c.JSON(201, gin.H{
		"message": "Chef created",
		"chef":    chef,
	})

}
func (h *Handler) getChef(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Get chef with id " + id,
	})
}

func (h *Handler) updateChef(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Update chef with id " + id,
	})
}
func (h *Handler) deleteChef(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Delete chef with id " + id,
	})
}

func (h *Handler) Login(c *gin.Context) {
	// Get the email and pass off the request body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Look up the requested user
	var chef domain.Chef
	h.DB.First(&chef, "email = ?", body.Email)

	if chef.ID == 0 {
		c.JSON(404, gin.H{
			"error": "Chef not found",
		})
		return
	}

	// Compare the provided password with the stored hashed password
	err := bcrypt.CompareHashAndPassword([]byte(chef.Password), []byte(body.Password))
	if err != nil {
		c.JSON(401, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// generate a JWT token if the password is correct
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": chef.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	//send it back
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func (h *Handler) Validate(c *gin.Context) {
	var user = c.MustGet("user").(domain.Chef)
	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
		"chef":    user,
	})
}
