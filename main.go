package main

import (
	"github.com/gin-gonic/gin"
	"recipea.com/m/Chef"
	"recipea.com/m/database"
	"recipea.com/m/domain"
)

func init() {
	database.InitDB()
	database.DB.AutoMigrate(&domain.Chef{})
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
