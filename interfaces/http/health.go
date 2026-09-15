package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"recipea.com/m/database"
)

func RegisterHealthRoutes(router *gin.Engine) {
	// Liveness check: confirms the process is running. This implementation also
	// checks the database, although liveness checks are usually dependency-free.
	router.GET("/health", func(c *gin.Context) {
		db, err := database.DB.DB()
		if err != nil || db == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "database unavailable"})
			return
		}
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Readiness check: confirms the service and its database are ready to receive traffic.
	router.GET("/ready", func(c *gin.Context) {
		if database.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		db, err := database.DB.DB()
		if err != nil || db == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
}
