package Chef

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"recipea.com/m/database"
	"recipea.com/m/middleware"
	domain "recipea.com/m/shared"
)

type Handler struct {
	DB      *gorm.DB
	Service *Service
}

var hmacSampleSecret = []byte("your-secret-key") // TODO: move this to an environment variable. Remove duplicate in requireAuth.go

func sanitizeChef(c domain.Chef) gin.H {
	return gin.H{
		"id":         c.ID,
		"first_name": c.FirstName,
		"last_name":  c.LastName,
		"email":      c.Email,
	}
}

func sanitizeChefs(chefs []domain.Chef) []gin.H {
	out := make([]gin.H, 0, len(chefs))
	for _, chef := range chefs {
		out = append(out, sanitizeChef(chef))
	}
	return out
}

func Routes(route *gin.Engine) {
	handler := &Handler{DB: database.DB, Service: NewService(database.DB)}
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
	chefs, err := h.Service.GetAllChefs()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch chefs"})
		return
	}
	c.IndentedJSON(http.StatusOK, sanitizeChefs(chefs))
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

	created, err := h.Service.CreateChef(domain.Chef{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create chef"})
		return
	}
	c.JSON(201, gin.H{"message": "Chef created", "chef": sanitizeChef(created)})

}
func (h *Handler) getChef(c *gin.Context) {
	id := c.Param("id")
	chef, err := h.Service.GetChefByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Chef not found"})
		return
	}
	c.JSON(200, gin.H{"chef": sanitizeChef(chef)})
}

func (h *Handler) updateChef(c *gin.Context) {
	id := c.Param("id")
	var body domain.Chef
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]any{}
	if body.FirstName != "" {
		updates["first_name"] = body.FirstName
	}
	if body.LastName != "" {
		updates["last_name"] = body.LastName
	}
	if body.Email != "" {
		updates["email"] = body.Email
	}
	if body.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		updates["password"] = string(hashedPassword)
	}

	if len(updates) == 0 {
		c.JSON(400, gin.H{"error": "No valid fields provided to update"})
		return
	}

	updatedChef, err := h.Service.UpdateChef(id, updates)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update chef"})
		return
	}

	c.JSON(200, gin.H{"message": "Updated chef with id " + id, "chef": sanitizeChef(updatedChef)})
}

func (h *Handler) deleteChef(c *gin.Context) {
	id := c.Param("id")
	if err := h.Service.DeleteChef(id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete chef"})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted chef with id " + id})
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

	chef, err := h.Service.FindChefByEmail(body.Email)
	if err != nil || chef.ID == 0 {
		c.JSON(404, gin.H{"error": "Chef not found"})
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(chef.Password), []byte(body.Password))
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
		"chef":    sanitizeChef(user),
	})
}
