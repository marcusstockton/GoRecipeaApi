package main

import (
	"github.com/gin-gonic/gin"
	appchef "recipea.com/m/application/chef"
	apprecipe "recipea.com/m/application/recipe"
	"recipea.com/m/database"
	"recipea.com/m/domain/chef"
	"recipea.com/m/domain/recipe"
	"recipea.com/m/infrastructure/auth"
	"recipea.com/m/infrastructure/persistence"
	httppkg "recipea.com/m/interfaces/http"
)

func init() {
	database.InitDB()
	database.DB.AutoMigrate(&chef.Chef{}, &recipe.Recipe{}, &recipe.Ingredient{}, &recipe.Step{})
}

func main() {
	router := gin.Default()
	jwtProvider := auth.NewJWTProvider(auth.DefaultSecret())
	chefRepo := persistence.NewChefRepository(database.DB)
	chefService := appchef.NewChefService(chefRepo, jwtProvider)

	recipeRepo := persistence.NewRecipeRepository(database.DB)
	recipeService := apprecipe.NewRecipeService(recipeRepo)

	httppkg.RegisterChefRoutes(router, chefService, jwtProvider, chefRepo)
	httppkg.RegisterRecipeRoutes(router, recipeService, jwtProvider, chefRepo)
	router.GET("/", homePage)

	router.Run("localhost:8080")
}

func homePage(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to the home page!"})
}
