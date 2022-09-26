package routes

import (
	"net/http"

	"github.com/Espigah/latency-aggregator/internal/domain/aggregator"
	"github.com/Espigah/latency-aggregator/internal/domain/appcontext"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MakeApplicationRoute creates a application route
func MakeApplicationRoute(r *gin.Engine, agg aggregator.Aggregator) {
	r.POST("/", func(c *gin.Context) {
		createMetric(c, agg)
	})
}

func createMetric(c *gin.Context, agg aggregator.Aggregator) {
	context := getContext(c)

	logger := context.Logger()
	logger.Info("Received message from integration")

	var histogramAggregator aggregator.MetricAggregatorDTO
	err := c.BindJSON(&histogramAggregator)

	if err != nil {
		logger.Error("Error binding JSON", zap.Error(err))
		respond(c, nil, err)
		return
	}
	logger.Debug("Received payload:", zap.Any("data", histogramAggregator))

	result, err := agg.CollectHistogram(histogramAggregator)

	respond(c, result, err)
}

//##########################################################################################

func getContext(c *gin.Context) appcontext.Context {
	return c.Value(string(appcontext.AppContextKey)).(appcontext.Context)
}

func respond(c *gin.Context, result interface{}, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "retryable": true})
		return
	}

	c.JSON(http.StatusOK, result)
}
