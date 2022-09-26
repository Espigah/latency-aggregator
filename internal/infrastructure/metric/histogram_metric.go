package metric

import (
	"github.com/Espigah/latency-aggregator/internal/domain/aggregator"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/metric/prometheus"
)

type histogramMetric struct {
	histogram           prometheus.HistogramVec
	metricAggregatorDTO aggregator.MetricAggregatorDTO
}

func (h *histogramMetric) Registry(duration float64) {
	h.histogram.Observe(duration, h.metricAggregatorDTO.LabelValues)
}

func NewHistogram(histogram prometheus.HistogramVec, metricAggregatorDTO aggregator.MetricAggregatorDTO) *histogramMetric {
	return &histogramMetric{
		histogram:           histogram,
		metricAggregatorDTO: metricAggregatorDTO,
	}
}
