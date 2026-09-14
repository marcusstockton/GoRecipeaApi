package http_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	apprecipe "recipea.com/m/application/recipe"
	domainchef "recipea.com/m/domain/chef"
	domainrecipe "recipea.com/m/domain/recipe"
	"recipea.com/m/infrastructure/auth"
	"recipea.com/m/infrastructure/persistence"
	httppkg "recipea.com/m/interfaces/http"
)

type recipeRouterTestHarness struct {
	router       *gin.Engine
	authProvider *auth.JWTProvider
	chefRepo     *persistence.ChefRepository
}

func setupRecipeRouterTest(t *testing.T) *recipeRouterTestHarness {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(
		sqlite.Open(filepath.Join(t.TempDir(), "recipe-test.db")),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying database connection: %v", err)
	}

	t.Cleanup(func() {
		sqlDB.Close()
	})

	if err := db.AutoMigrate(
		&domainchef.Chef{},
		&domainrecipe.Recipe{},
		&domainrecipe.Ingredient{},
		&domainrecipe.Step{},
		&domainrecipe.RecipeLike{},
		&domainrecipe.RecipeComment{},
	); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	chefRepo := persistence.NewChefRepository(db)
	recipeRepo := persistence.NewRecipeRepository(db)
	authProvider := auth.NewJWTProvider("test-secret")
	service := apprecipe.NewRecipeService(recipeRepo)

	router := gin.New()
	httppkg.RegisterRecipeRoutes(router, service, authProvider, chefRepo)

	return &recipeRouterTestHarness{
		router:       router,
		authProvider: authProvider,
		chefRepo:     chefRepo,
	}
}

func (h *recipeRouterTestHarness) createChef(t *testing.T, firstName, lastName, email string) uint {
	t.Helper()
	chef, err := domainchef.NewChef(firstName, lastName, email, "password123")
	if err != nil {
		t.Fatalf("failed to create chef: %v", err)
	}
	if err := h.chefRepo.Save(chef); err != nil {
		t.Fatalf("failed to persist chef: %v", err)
	}
	return chef.ID
}

func (h *recipeRouterTestHarness) tokenFor(t *testing.T, chefID uint) string {
	t.Helper()
	token, err := h.authProvider.CreateToken(chefID)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}
	return token
}

func doRecipeRequest(t *testing.T, router *gin.Engine, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	router.ServeHTTP(w, req)
	return w
}

func TestRecipeRoutes_EnforceAuthAndLikeCommentOwnership(t *testing.T) {
	h := setupRecipeRouterTest(t)
	ownerID := h.createChef(t, "Mina", "Bake", "mina@example.com")
	viewerID := h.createChef(t, "Jules", "Cook", "jules@example.com")
	ownerToken := h.tokenFor(t, ownerID)
	viewerToken := h.tokenFor(t, viewerID)

	createBody := `{"title":"Curry","description":"Spicy and warm","ingredients":[{"quantity":"1","unit":"tbsp","name":"curry powder"}],"steps":[{"order":1,"action":"Cook"}]}`
	w := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/", createBody, ownerToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected create recipe to return 201, got %d body=%s", w.Code, w.Body.String())
	}

	unauthenticatedLike := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/1/like", "", "")
	if unauthenticatedLike.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated like to return 401, got %d body=%s", unauthenticatedLike.Code, unauthenticatedLike.Body.String())
	}

	ownerLike := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/1/like", "", ownerToken)
	if ownerLike.Code != http.StatusBadRequest {
		t.Fatalf("expected owner self-like to return 400, got %d body=%s", ownerLike.Code, ownerLike.Body.String())
	}

	viewerLike := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/1/like", "", viewerToken)
	if viewerLike.Code != http.StatusCreated {
		t.Fatalf("expected viewer like to return 201, got %d body=%s", viewerLike.Code, viewerLike.Body.String())
	}

	duplicateLike := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/1/like", "", viewerToken)
	if duplicateLike.Code != http.StatusConflict {
		t.Fatalf("expected duplicate like to return 409, got %d body=%s", duplicateLike.Code, duplicateLike.Body.String())
	}

	unauthorizedRemoveLike := doRecipeRequest(t, h.router, http.MethodDelete, "/recipe/1/like", "", ownerToken)
	if unauthorizedRemoveLike.Code != http.StatusForbidden {
		t.Fatalf("expected removing another chef's like to return 403, got %d body=%s", unauthorizedRemoveLike.Code, unauthorizedRemoveLike.Body.String())
	}

	removeLike := doRecipeRequest(t, h.router, http.MethodDelete, "/recipe/1/like", "", viewerToken)
	if removeLike.Code != http.StatusOK {
		t.Fatalf("expected removing own like to return 200, got %d body=%s", removeLike.Code, removeLike.Body.String())
	}

	commentBody := `{"content":"Great recipe"}`
	ownerComment := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/1/comment", commentBody, ownerToken)
	if ownerComment.Code != http.StatusCreated {
		t.Fatalf("expected owner comment to return 201, got %d body=%s", ownerComment.Code, ownerComment.Body.String())
	}

	comments := doRecipeRequest(t, h.router, http.MethodGet, "/recipe/1/comments", "", "")
	if comments.Code != http.StatusOK {
		t.Fatalf("expected comments listing to return 200, got %d body=%s", comments.Code, comments.Body.String())
	}

	forbiddenRemoveComment := doRecipeRequest(t, h.router, http.MethodDelete, "/recipe/comment/1", "", viewerToken)
	if forbiddenRemoveComment.Code != http.StatusForbidden {
		t.Fatalf("expected removing another chef's comment to return 403, got %d body=%s", forbiddenRemoveComment.Code, forbiddenRemoveComment.Body.String())
	}

	removeComment := doRecipeRequest(t, h.router, http.MethodDelete, "/recipe/comment/1", "", ownerToken)
	if removeComment.Code != http.StatusOK {
		t.Fatalf("expected removing own comment to return 200, got %d body=%s", removeComment.Code, removeComment.Body.String())
	}
}

func TestRecipeRoutes_RejectReplyToCommentFromAnotherRecipe(t *testing.T) {
	h := setupRecipeRouterTest(t)
	chefID := h.createChef(t, "Rae", "Cook", "rae@example.com")
	token := h.tokenFor(t, chefID)

	createRecipeBody := `{"title":"Pie","description":"Sweet","ingredients":[{"quantity":"1","unit":"slice","name":"pie"}],"steps":[{"order":1,"action":"Bake"}]}`
	firstRecipe := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/", createRecipeBody, token)
	if firstRecipe.Code != http.StatusCreated {
		t.Fatalf("expected first recipe creation to return 201, got %d body=%s", firstRecipe.Code, firstRecipe.Body.String())
	}

	secondRecipeBody := `{"title":"Cake","description":"Sweet","ingredients":[{"quantity":"1","unit":"slice","name":"cake"}],"steps":[{"order":1,"action":"Bake"}]}`
	secondRecipe := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/", secondRecipeBody, token)
	if secondRecipe.Code != http.StatusCreated {
		t.Fatalf("expected second recipe creation to return 201, got %d body=%s", secondRecipe.Code, secondRecipe.Body.String())
	}

	commentBody := `{"content":"First recipe comment"}`
	createdComment := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/1/comment", commentBody, token)
	if createdComment.Code != http.StatusCreated {
		t.Fatalf("expected comment creation to return 201, got %d body=%s", createdComment.Code, createdComment.Body.String())
	}

	replyBody := `{"content":"Replying to the wrong recipe","parent_id":1}`
	wrongRecipeReply := doRecipeRequest(t, h.router, http.MethodPost, "/recipe/2/comment", replyBody, token)
	if wrongRecipeReply.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid cross-recipe reply to return 400, got %d body=%s", wrongRecipeReply.Code, wrongRecipeReply.Body.String())
	}
}
