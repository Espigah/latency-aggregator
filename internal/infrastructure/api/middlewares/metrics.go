package middlewares

import (
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware is a middleware that records metrics
func MetricsMiddleware() gin.HandlerFunc {
	// prometheusService, err := metric.NewPrometheusService()
	// if err != nil {
	// 	panic(err)
	// }
	return func(context *gin.Context) {
		// appMetric := metric.NewHTTP(context.FullPath(), context.Request.Method)

		// appMetric.Started()

		context.Next()

		// appMetric.Finished()

		// appMetric.StatusCode = strconv.Itoa(context.Writer.Status())

		// prometheusService.SaveHTTP(appMetric)
	}
}
