package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	appchef "recipea.com/m/application/chef"
	apprecipe "recipea.com/m/application/recipe"
	"recipea.com/m/config"
	"recipea.com/m/database"
	"recipea.com/m/domain/chef"
	"recipea.com/m/domain/recipe"
	"recipea.com/m/infrastructure/auth"
	"recipea.com/m/infrastructure/logging"
	"recipea.com/m/infrastructure/persistence"
	httppkg "recipea.com/m/interfaces/http"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)
	logging.Init(cfg.LogLevel)

	slog.Info("starting application", "host", cfg.Host, "port", cfg.Port, "database_path", cfg.DatabasePath)
	slog.Info("resolved database path", "path", cfg.DatabasePath)

	database.InitDB(cfg.DatabasePath)
	database.DB.AutoMigrate(&chef.Chef{}, &recipe.Recipe{}, &recipe.Ingredient{}, &recipe.Step{}, &recipe.RecipeLike{}, &recipe.RecipeComment{})
	if err := database.Seed(); err != nil {
		slog.Error("database seed failed", "error", err)
	}

	router := gin.New()
	router.Use(httppkg.RequestIDMiddleware())
	router.Use(httppkg.LoggerMiddleware())
	router.Use(httppkg.ErrorHandler())

	jwtProvider := auth.NewJWTProvider(cfg.JWTSecret)
	chefRepo := persistence.NewChefRepository(database.DB)
	chefService := appchef.NewChefService(chefRepo, jwtProvider)

	recipeRepo := persistence.NewRecipeRepository(database.DB)
	recipeService := apprecipe.NewRecipeService(recipeRepo)

	httppkg.RegisterHealthRoutes(router)
	httppkg.RegisterChefRoutes(router, chefService, jwtProvider, chefRepo)
	httppkg.RegisterRecipeRoutes(router, recipeService, jwtProvider, chefRepo)
	router.GET("/", homePage)

	if err := router.Run(cfg.Addr()); err != nil {
		slog.Error("server shutdown", "error", err)
	}
}

func homePage(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to the home page!"})
}
