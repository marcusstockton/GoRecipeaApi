package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apprecipe "recipea.com/m/application/recipe"
	domainchef "recipea.com/m/domain/chef"
	domainrecipe "recipea.com/m/domain/recipe"
	"recipea.com/m/infrastructure/auth"
)

type RecipeHandler struct {
	Service *apprecipe.RecipeService
}

type ingredientItem struct {
	Quantity    string `json:"quantity" binding:"required"`
	Unit        string `json:"unit" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Preparation string `json:"preparation,omitempty"`
}

type stepItem struct {
	Order  int    `json:"order" binding:"required"`
	Action string `json:"action" binding:"required"`
}

type createRecipeRequest struct {
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description"`
	Ingredients []ingredientItem `json:"ingredients" binding:"required,dive,required"`
	Steps       []stepItem       `json:"steps" binding:"required,dive,required"`
}

type updateRecipeRequest struct {
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description"`
	Ingredients []ingredientItem `json:"ingredients" binding:"required,dive,required"`
	Steps       []stepItem       `json:"steps" binding:"required,dive,required"`
}

type addCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

func RegisterRecipeRoutes(router *gin.Engine, service *apprecipe.RecipeService, authProvider *auth.JWTProvider, chefRepo domainchef.Repository) {
	h := &RecipeHandler{Service: service}
	rg := router.Group("/recipe")
	rg.GET("/", h.ListRecipes)
	rg.GET("/:id", h.GetRecipe)
	rg.GET("/:id/comments", h.ListComments)
	authGroup := rg.Group("/")
	authGroup.Use(RequireAuth(authProvider, chefRepo))
	authGroup.POST("/", h.CreateRecipe)
	authGroup.PUT("/:id", h.UpdateRecipe)
	authGroup.DELETE("/:id", h.DeleteRecipe)
	authGroup.POST("/:id/like", h.LikeRecipe)
	authGroup.DELETE("/:id/like", h.RemoveLike)
	authGroup.POST("/:id/comment", h.AddComment)
	authGroup.DELETE("/comment/:commentID", h.RemoveComment)
}

func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	chefObj, exists := c.Get("chef")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "chef context missing"})
		return
	}
	chef, ok := chefObj.(*domainchef.Chef)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid chef context"})
		return
	}
	var req createRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ingredients := make([]domainrecipe.Ingredient, 0, len(req.Ingredients))
	for _, i := range req.Ingredients {
		ingredients = append(ingredients, domainrecipe.Ingredient{Quantity: i.Quantity, Unit: i.Unit, Name: i.Name, Preparation: i.Preparation})
	}
	steps := make([]domainrecipe.Step, 0, len(req.Steps))
	for _, s := range req.Steps {
		steps = append(steps, domainrecipe.Step{Order: s.Order, Action: s.Action})
	}
	out, err := h.Service.CreateRecipe(apprecipe.CreateRecipeRequest{Title: req.Title, Description: req.Description, OwnedBy: chef.ID, Ingredients: ingredients, Steps: steps})
	if err != nil {
		status := http.StatusBadRequest
		if err == domainrecipe.ErrNotAuthorized {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"recipe": out})
}

func (h *RecipeHandler) ListRecipes(c *gin.Context) {
	out, err := h.Service.ListRecipes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	out, err := h.Service.GetRecipe(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recipe": out})
}

func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	chefObj, exists := c.Get("chef")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "chef context missing"})
		return
	}
	chef, ok := chefObj.(*domainchef.Chef)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid chef context"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ingredients := make([]domainrecipe.Ingredient, 0, len(req.Ingredients))
	for _, i := range req.Ingredients {
		ingredients = append(ingredients, domainrecipe.Ingredient{Quantity: i.Quantity, Unit: i.Unit, Name: i.Name, Preparation: i.Preparation})
	}
	steps := make([]domainrecipe.Step, 0, len(req.Steps))
	for _, s := range req.Steps {
		steps = append(steps, domainrecipe.Step{Order: s.Order, Action: s.Action})
	}
	out, err := h.Service.UpdateRecipe(apprecipe.UpdateRecipeRequest{ID: uint(id), UpdatedBy: chef.ID, Title: req.Title, Description: req.Description, Ingredients: ingredients, Steps: steps})
	if err != nil {
		status := http.StatusBadRequest
		if err == domainrecipe.ErrNotAuthorized {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recipe": out})
}

func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	chefObj, exists := c.Get("chef")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "chef context missing"})
		return
	}
	chef, ok := chefObj.(*domainchef.Chef)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid chef context"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Service.DeleteRecipe(uint(id), chef.ID); err != nil {
		status := http.StatusInternalServerError
		if err == domainrecipe.ErrNotAuthorized {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *RecipeHandler) LikeRecipe(c *gin.Context) {
	chefObj, exists := c.Get("chef")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "chef context missing"})
		return
	}
	chef, ok := chefObj.(*domainchef.Chef)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid chef context"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Service.AddLike(uint(id), chef.ID); err != nil {
		status := http.StatusBadRequest
		switch err {
		case domainrecipe.ErrDuplicateLike:
			status = http.StatusConflict
		case domainrecipe.ErrNotFound:
			status = http.StatusNotFound
		case domainrecipe.ErrInvalidInput, domainrecipe.ErrSelfLikeNotAllowed:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "liked"})
}

func (h *RecipeHandler) RemoveLike(c *gin.Context) {
	chefObj, exists := c.Get("chef")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "chef context missing"})
		return
	}
	chef, ok := chefObj.(*domainchef.Chef)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid chef context"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Service.RemoveLike(uint(id), chef.ID); err != nil {
		status := http.StatusInternalServerError
		switch err {
		case domainrecipe.ErrNotAuthorized:
			status = http.StatusForbidden
		case domainrecipe.ErrNotFound:
			status = http.StatusNotFound
		case domainrecipe.ErrInvalidInput:
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "like removed"})
}

func (h *RecipeHandler) AddComment(c *gin.Context) {
	chefObj, exists := c.Get("chef")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "chef context missing"})
		return
	}
	chef, ok := chefObj.(*domainchef.Chef)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid chef context"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req addCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.AddComment(uint(id), chef.ID, req.Content, req.ParentID); err != nil {
		status := http.StatusBadRequest
		switch err {
		case domainrecipe.ErrEmptyComment:
			status = http.StatusBadRequest
		case domainrecipe.ErrNotFound:
			status = http.StatusNotFound
		case domainrecipe.ErrInvalidInput:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "comment added"})
}

func (h *RecipeHandler) RemoveComment(c *gin.Context) {
	chefObj, exists := c.Get("chef")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "chef context missing"})
		return
	}
	chef, ok := chefObj.(*domainchef.Chef)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid chef context"})
		return
	}
	commentID, err := strconv.ParseUint(c.Param("commentID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}
	if err := h.Service.RemoveComment(uint(commentID), chef.ID); err != nil {
		status := http.StatusInternalServerError
		switch err {
		case domainrecipe.ErrNotAuthorized:
			status = http.StatusForbidden
		case domainrecipe.ErrNotFound:
			status = http.StatusNotFound
		case domainrecipe.ErrInvalidInput:
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "comment removed"})
}

func (h *RecipeHandler) ListComments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	comments, err := h.Service.ListComments(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
