package api

import (
	"github.com/gin-gonic/gin"

	"github.com/Espigah/latency-aggregator/internal/domain/aggregator"
	"github.com/Espigah/latency-aggregator/internal/environment"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/api/middlewares"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/api/routes"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/logger/logwrapper"
	"go.uber.org/zap"
)

// Input is the input for the API server
type Input struct {
	Logger     logwrapper.LoggerWrapper
	Aggregator aggregator.Aggregator
}

// Start creates a new API server
func Start(input Input) {
	r := gin.New()
	env := environment.GetInstance()

	logger := input.Logger

	logger.Info("Starting latency-aggregator API")

	applicationPort := resolvePort()
	r.Use(middlewares.MetricsMiddleware())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.ContextMiddleware())
	r.Use(middlewares.TraceMiddleware())
	r.Use(middlewares.Logger(logger))

	if !env.IsDevelopment() {
		r.Use(middlewares.Recovery(true))
	}

	routes.MakeHealthRoute(r)
	routes.MakeMetricRoute(r)
	routes.MakeApplicationRoute(r, input.Aggregator)

	if err := r.Run(applicationPort); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}

func resolvePort() string {
	const CHAR string = ":"
	env := environment.GetInstance()
	port := env.APP_PORT
	fisrtChar := port[:1]
	if fisrtChar != CHAR {
		port = CHAR + port
	}
	return port
}
