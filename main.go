package main

import (
	"github.com/gin-gonic/gin"
	Chef "recipea.com/m/chef" // alias with capital C
	"recipea.com/m/database"
	"recipea.com/m/models"
)

func init() {
	database.InitDB()
	database.DB.AutoMigrate(&models.Chef{})
}

func main() {

	router := gin.Default()
	Chef.Routes(router)
	router.GET("/", homePage)

	router.Run("localhost:8080")
}

func homePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the home page!",
	})
}
