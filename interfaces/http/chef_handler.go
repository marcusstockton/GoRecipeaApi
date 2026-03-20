package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	appchef "recipea.com/m/application/chef"
	"recipea.com/m/domain/chef"
	"recipea.com/m/infrastructure/auth"
)

type ChefHandler struct {
	Service *appchef.ChefService
	Auth    *auth.JWTProvider
}

type createChefRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type updateChefRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func RegisterChefRoutes(router *gin.Engine, service *appchef.ChefService, authProvider *auth.JWTProvider, repo chef.Repository) {
	h := &ChefHandler{Service: service, Auth: authProvider}
	chefRG := router.Group("/chef")
	chefRG.GET("/", h.GetChefs)
	chefRG.GET("/:id", h.GetChef)
	chefRG.POST("/", h.CreateChef)
	chefRG.PUT("/:id", h.UpdateChef)
	chefRG.POST("/login", h.Login)
	chefRG.GET("/validate", RequireAuth(authProvider, repo), h.Validate)
}

func (h *ChefHandler) GetChefs(c *gin.Context) {
	chefs, err := h.Service.ListChefs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chefs)
}

func (h *ChefHandler) CreateChef(c *gin.Context) {
	var req createChefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.Service.CreateChef(appchef.CreateChefRequest{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"chef": out})
}

func (h *ChefHandler) GetChef(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	out, err := h.Service.GetChef(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"chef": out})
}

func (h *ChefHandler) UpdateChef(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateChefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.Service.UpdateChef(appchef.UpdateChefRequest{ID: uint(id), FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"chef": out})
}

func (h *ChefHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, user, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "chef": user})
}

func (h *ChefHandler) Validate(c *gin.Context) {
	chefObj, exists := c.Get("chef")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing chef context"})
		return
	}
	chef := chefObj.(*chef.Chef)
	c.JSON(http.StatusOK, gin.H{"chef": gin.H{"id": chef.ID, "email": chef.Email, "first_name": chef.FirstName, "last_name": chef.LastName}})
}
