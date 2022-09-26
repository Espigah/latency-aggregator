package internal_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Espigah/latency-aggregator/internal/domain/aggregator"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/api"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/collector"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/database"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/logger"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/logger/logwrapper"
	latencyaggregator "github.com/Espigah/latency-aggregator/pkg/prometheus/latency-aggregator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	const url = "http://localhost:7070"
	zaplogger, dispose := logger.New()
	defer dispose()

	logger := logwrapper.New(&logwrapper.Zap{Logger: *zaplogger})

	input := aggregator.Input{
		Repository: database.NewMemoryDatabase(),
		Collector:  collector.New(),
	}

	aggregator, err := aggregator.New(input)

	if err != nil {
		logger.Error("failed to setup Aggregator")
		panic(err)
	}

	apiInput := api.Input{
		Logger:     logger,
		Aggregator: aggregator,
	}

	go api.Start(apiInput)

	t.Run("scenario: creating a simple metric", func(t *testing.T) {
		var (
			name          = "my_histogram"
			description   = "it's just a test"
			wantHistogram = `my_histogram_count 1`
			wantGouge     = "my_histogram_total"
		)

		latencyaggregator.New().
			ID(uuid.New().String()).
			URL(url).
			LifeSpan(0).
			Name(name).
			Help(description).
			Push()
		time.Sleep(1 * time.Second)
		assertHistogramSum(t, url, wantHistogram)
		assertHistogramSum(t, url, wantGouge)
	})

	t.Run("scenario: creating metric with labels", func(t *testing.T) {
		var (
			name          = "my_histogram_with_labels"
			description   = "it's just a test"
			wantHistogram = `my_histogram_with_labels_count{method="GET",status="200"} 1`
			wantGouge     = `my_histogram_with_labels_total{method="GET",status="200"}`
		)

		labels := map[string]string{
			"status": "200",
			"method": "GET",
		}
		latencyaggregator.New().
			ID(uuid.New().String()).
			URL(url).
			LifeSpan(0).
			Name(name).
			LabelValues(labels).
			Help(description).
			Push()
		time.Sleep(1 * time.Second)
		assertHistogramSum(t, url, wantHistogram)
		assertHistogramSum(t, url, wantGouge)
	})

	t.Run("scenario: creating 2 metric with labels different labels and same stage", func(t *testing.T) {
		var (
			name          = "my_histogram_different_labels"
			description   = "it's just a test"
			wantHistogram = `my_histogram_different_labels_count{method="GET",status="200"} 1`
			wantGouge     = `my_histogram_different_labels_total{method="GET",status="200"}`
		)

		labels := map[string]string{
			"status": "200",
			"method": "GET",
		}
		latencyaggregator.New().
			ID(uuid.New().String()).
			URL(url).
			LifeSpan(2).
			Name(name).
			LabelValues(labels).
			Help(description).
			Push()
		time.Sleep(1 * time.Second)
		labels2 := map[string]string{
			"status": "300",
			"method": "test",
		}
		latencyaggregator.New().
			ID(uuid.New().String()).
			URL(url).
			LifeSpan(0).
			Name(name).
			LabelValues(labels2).
			Help(description).
			Push()
		time.Sleep(1 * time.Second)
		assertHistogramSum(t, url, wantHistogram)
		assertHistogramSum(t, url, wantGouge)
	})
}

func assertHistogramSum(t *testing.T, url, want string) {
	r, err := http.Get(fmt.Sprintf("%s%s", url, "/metrics"))
	assert.NoError(t, err)
	b, err := io.ReadAll(r.Body)
	assert.NoError(t, err)
	body := string(b)
	assert.Contains(t, body, want)
}
