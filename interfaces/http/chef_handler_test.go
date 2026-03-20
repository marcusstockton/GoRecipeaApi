package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	appchef "recipea.com/m/application/chef"
	"recipea.com/m/domain/chef"
	"recipea.com/m/infrastructure/auth"
	"recipea.com/m/infrastructure/persistence"
	httppkg "recipea.com/m/interfaces/http"
)

type createChefResp struct {
	Chef struct {
		ID        uint   `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	} `json:"chef"`
}

type loginResp struct {
	Token string `json:"token"`
}

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&chef.Chef{}); err != nil {
		t.Fatalf("failed migrate: %v", err)
	}

	repo := persistence.NewChefRepository(db)
	jwtProvider := auth.NewJWTProvider("test-secret")
	service := appchef.NewChefService(repo, jwtProvider)

	router := gin.Default()
	httppkg.RegisterChefRoutes(router, service, jwtProvider, repo)
	return router
}

func TestChefEndpoints_CreateLoginValidate(t *testing.T) {
	router := setupTestRouter(t)

	// Create chef
	createReq := `{"first_name":"Mina","last_name":"Bake","email":"mina@recipe.com","password":"password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/chef/", strings.NewReader(createReq))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", w.Code, w.Body.String())
	}

	var created createChefResp
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed unmarshal create response: %v", err)
	}
	if created.Chef.Email != "mina@recipe.com" {
		t.Fatalf("expected chef email mina@recipe.com, got %s", created.Chef.Email)
	}

	// Login
	loginBody := `{"email":"mina@recipe.com","password":"password123"}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/chef/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 login, got %d body=%s", w.Code, w.Body.String())
	}

	var login loginResp
	if err := json.Unmarshal(w.Body.Bytes(), &login); err != nil {
		t.Fatalf("failed unmarshal login response: %v", err)
	}
	if login.Token == "" {
		t.Fatal("expected token in login response")
	}

	// Validate
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/chef/validate", nil)
	req.Header.Set("Authorization", "Bearer "+login.Token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 validate, got %d body=%s", w.Code, w.Body.String())
	}
}
