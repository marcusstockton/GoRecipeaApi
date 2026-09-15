package main

import (
	"github.com/gin-gonic/gin"
	appchef "recipea.com/m/application/chef"
	apprecipe "recipea.com/m/application/recipe"
	"recipea.com/m/config"
	"recipea.com/m/database"
	"recipea.com/m/domain/chef"
	"recipea.com/m/domain/recipe"
	"recipea.com/m/infrastructure/auth"
	"recipea.com/m/infrastructure/persistence"
	httppkg "recipea.com/m/interfaces/http"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	database.InitDB(cfg.DatabasePath)
	database.DB.AutoMigrate(&chef.Chef{}, &recipe.Recipe{}, &recipe.Ingredient{}, &recipe.Step{}, &recipe.RecipeLike{}, &recipe.RecipeComment{})
	_ = database.Seed()

	router := gin.Default()
	jwtProvider := auth.NewJWTProvider(cfg.JWTSecret)
	chefRepo := persistence.NewChefRepository(database.DB)
	chefService := appchef.NewChefService(chefRepo, jwtProvider)

	recipeRepo := persistence.NewRecipeRepository(database.DB)
	recipeService := apprecipe.NewRecipeService(recipeRepo)

	httppkg.RegisterChefRoutes(router, chefService, jwtProvider, chefRepo)
	httppkg.RegisterRecipeRoutes(router, recipeService, jwtProvider, chefRepo)
	router.GET("/", homePage)

	router.Run(cfg.Addr())
}

func homePage(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to the home page!"})
}
