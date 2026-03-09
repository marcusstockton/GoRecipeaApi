package chef

import "github.com/gin-gonic/gin"

func Routes(route *gin.Engine) {
	chef := route.Group("/chef")
	{
		chef.GET("/", getChefs)
		chef.GET("/:id", getChef)
		chef.POST("/", createChef)
		chef.PUT("/:id", updateChef)
		chef.DELETE("/:id", deleteChef)
	}
}

func getChefs(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get chefs",
	})
}

func createChef(c *gin.Context) {
	c.JSON(201, gin.H{
		"message": "Create chef called",
	})

}
func getChef(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Get chef with id " + id,
	})
}

func updateChef(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Update chef with id " + id,
	})
}
func deleteChef(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Delete chef with id " + id,
	})
}
