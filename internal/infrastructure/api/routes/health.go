package routes

import (
	"net/http"

	"github.com/Espigah/latency-aggregator/internal/environment"
	"github.com/gin-gonic/gin"
)

// MakeHealthRoute creates a health route
func MakeHealthRoute(r *gin.Engine) {
	version := environment.GetInstance().APP_VERSION
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "version": version})
	})
}
