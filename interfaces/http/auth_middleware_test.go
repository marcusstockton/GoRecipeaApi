package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"recipea.com/m/domain/chef"
	"recipea.com/m/infrastructure/auth"
)

type fakeChefRepository struct {
	chef     *chef.Chef
	findErr  error
	findID   uint
	findCall int
}

func (r *fakeChefRepository) Save(c *chef.Chef) error {
	return nil
}

func (r *fakeChefRepository) FindByID(id uint) (*chef.Chef, error) {
	r.findCall++
	r.findID = id

	if r.findErr != nil {
		return nil, r.findErr
	}

	return r.chef, nil

}

func (r *fakeChefRepository) FindByEmail(email string) (*chef.Chef, error) {
	return nil, nil
}

func (r *fakeChefRepository) Delete(id uint) error {
	return nil
}

func (r *fakeChefRepository) List() ([]chef.Chef, error) {
	return nil, nil
}

func setupAuthMiddlewareTest(t *testing.T, repo chef.Repository) (*gin.Engine, *auth.JWTProvider) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	provider := auth.NewJWTProvider("test-secret")

	router := gin.New()
	router.Use(RequireAuth(provider, repo))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router, provider

}

func doAuthMiddlewareRequest(
	t *testing.T,
	router *gin.Engine,
	token string,
) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder

}

func TestRequireAuth_NoAuthorizationHeader(t *testing.T) {
	repo := &fakeChefRepository{}
	router, _ := setupAuthMiddlewareTest(t, repo)

	response := doAuthMiddlewareRequest(t, router, "")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status 401, got %d body=%s",
			response.Code,
			response.Body.String(),
		)
	}

	if repo.findCall != 0 {
		t.Fatalf("expected repository not to be called, got %d calls", repo.findCall)
	}

}

func TestRequireAuth_InvalidToken(t *testing.T) {
	repo := &fakeChefRepository{}
	router, _ := setupAuthMiddlewareTest(t, repo)

	response := doAuthMiddlewareRequest(t, router, "this-is-not-a-valid-token")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status 401, got %d body=%s",
			response.Code,
			response.Body.String(),
		)
	}

	if repo.findCall != 0 {
		t.Fatalf("expected repository not to be called, got %d calls", repo.findCall)
	}

}

func TestRequireAuth_UserNotFound(t *testing.T) {
	repo := &fakeChefRepository{
		findErr: errors.New("chef not found"),
	}

	router, provider := setupAuthMiddlewareTest(t, repo)

	token, err := provider.CreateToken(123)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	response := doAuthMiddlewareRequest(t, router, token)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status 401, got %d body=%s",
			response.Code,
			response.Body.String(),
		)
	}

	if repo.findCall != 1 {
		t.Fatalf("expected repository to be called once, got %d calls", repo.findCall)
	}

	if repo.findID != 123 {
		t.Fatalf("expected repository to look up chef 123, got %d", repo.findID)
	}

}

func TestRequireAuth_ValidToken(t *testing.T) {
	testChef, err := chef.NewChef(
		"Mina",
		"Bake",
		"mina@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("failed to create chef: %v", err)
	}

	repo := &fakeChefRepository{
		chef: testChef,
	}

	router, provider := setupAuthMiddlewareTest(t, repo)

	token, err := provider.CreateToken(testChef.ID)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	response := doAuthMiddlewareRequest(t, router, token)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d body=%s",
			response.Code,
			response.Body.String(),
		)
	}

	if repo.findCall != 1 {
		t.Fatalf("expected repository to be called once, got %d calls", repo.findCall)
	}

	if repo.findID != testChef.ID {
		t.Fatalf(
			"expected repository to look up chef %d, got %d",
			testChef.ID,
			repo.findID,
		)
	}
}

func TestRequireAuth_ValidTokenUsesChefFromRepository(t *testing.T) {
	testChef, err := chef.NewChef(
		"Jules",
		"Cook",
		"jules@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("failed to create chef: %v", err)
	}

	repo := &fakeChefRepository{
		chef: testChef,
	}

	gin.SetMode(gin.TestMode)
	provider := auth.NewJWTProvider("test-secret")

	router := gin.New()
	router.Use(RequireAuth(provider, repo))

	router.GET("/protected", func(c *gin.Context) {
		chefObj, exists := c.Get("chef")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "chef not found in context",
			})
			return
		}

		contextChef, ok := chefObj.(*chef.Chef)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid chef in context",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"chef_id": contextChef.ID,
		})
	})

	token, err := provider.CreateToken(testChef.ID)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	response := doAuthMiddlewareRequest(t, router, token)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d body=%s",
			response.Code,
			response.Body.String(),
		)
	}

	expectedBody := `{"chef_id":` + strconv.FormatUint(uint64(testChef.ID), 10) + `}`

	if response.Body.String() != expectedBody {
		t.Fatalf(
			"expected body %s, got %s",
			expectedBody,
			response.Body.String(),
		)
	}
}
