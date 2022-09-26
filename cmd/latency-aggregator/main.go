package main

import (
	"github.com/Espigah/latency-aggregator/internal/domain/aggregator"
	"github.com/Espigah/latency-aggregator/internal/environment"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/api"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/collector"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/database"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/logger"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/logger/logwrapper"
	"go.uber.org/zap"
)

func main() {

	logger, dispose := setupLogger()
	defer dispose()
	logger.Info("Starting Histogram Aggregator ")

	aggregator, err := setupAggregator()

	if err != nil {
		logger.Error("failed to setup Aggregator", zap.Error(err))
		panic(err)
	}

	setupApi(logger, aggregator)
}

func setupLogger() (logwrapper.LoggerWrapper, func()) {
	env := environment.GetInstance()
	zaplogger, dispose := logger.New()

	logger := logwrapper.New(&logwrapper.Zap{Logger: *zaplogger}).SetVersion(env.APP_VERSION)

	logger.Info("env",
		zap.String("LOG_LEVEL", env.LOG_LEVEL),
		zap.String("APP_PORT", env.APP_PORT),
		zap.String("ENVIRONMENT", env.ENVIRONMENT),
		zap.String("VERSION", env.APP_VERSION),
	)

	return logger, dispose
}

func setupAggregator() (aggregator.Aggregator, error) {

	repository, err := setupRepository()
	if err != nil {
		return nil, err
	}

	collector, err := setupCollector()
	if err != nil {
		return nil, err
	}

	input := aggregator.Input{
		Repository: repository,
		Collector:  collector,
	}

	aggregator, err := aggregator.New(input)

	if err != nil {
		return nil, err
	}

	return aggregator, nil
}

func setupRepository() (aggregator.Repository, error) {
	return database.NewMemoryDatabase(), nil

}

func setupApi(logger logwrapper.LoggerWrapper, aggregator aggregator.Aggregator) {
	input := api.Input{
		Logger:     logger,
		Aggregator: aggregator,
	}

	api.Start(input)
}

func setupCollector() (aggregator.Collector, error) {
	return collector.New(), nil
}
