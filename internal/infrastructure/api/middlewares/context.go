package middlewares

import (
	"context"

	"github.com/Espigah/latency-aggregator/internal/domain/appcontext"
	"github.com/gin-gonic/gin"
)

// ContextMiddleware create appcontext
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		ctx := context.Background()
		appContext := appcontext.New(ctx, c)

		c.Set(string(appcontext.AppContextKey), appContext)

		c.Next()

		appContext.Done()
		c.Set(string(appcontext.AppContextKey), nil)
	}
}
