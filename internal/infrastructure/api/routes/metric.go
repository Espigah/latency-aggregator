package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MakeMetricRoute creates the routes for the metrics
func MakeMetricRoute(r *gin.Engine) {
	r.GET("/metrics", func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})
}
