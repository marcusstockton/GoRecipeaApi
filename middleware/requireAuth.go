package middleware

import (
	"github.com/gin-gonic/gin"
	"recipea.com/m/domain/chef"
	"recipea.com/m/infrastructure/auth"
	httpmiddleware "recipea.com/m/interfaces/http"
)

// Deprecated: use interfaces/http.RequireAuth with an injected repository.
// This adapter keeps legacy callers working without reintroducing direct DB access into the HTTP layer.
func RequireAuth(repo chef.Repository) gin.HandlerFunc {
	provider := auth.NewJWTProvider(auth.DefaultSecret())
	return httpmiddleware.RequireAuth(provider, repo)
}
