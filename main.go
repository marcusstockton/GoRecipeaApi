package main

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

func main() {

	router := gin.Default()
	router.GET("/", homePage)

	//db.AutoMigrate(&album{})

	router.Run("localhost:8080")
}

func homePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the home page!",
	})
}
