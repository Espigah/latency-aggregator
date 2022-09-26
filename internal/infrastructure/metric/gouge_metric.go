package metric

import (
	"github.com/Espigah/latency-aggregator/internal/domain/aggregator"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/metric/prometheus"
)

type gougeMetric struct {
	gouge               prometheus.GaugeVec
	metricAggregatorDTO aggregator.MetricAggregatorDTO
}

func (h *gougeMetric) Registry(duration float64) {
	h.gouge.Observe(duration, h.metricAggregatorDTO.LabelValues)
}

func NewGouge(gouge prometheus.GaugeVec, metricAggregatorDTO aggregator.MetricAggregatorDTO) *gougeMetric {
	return &gougeMetric{
		gouge:               gouge,
		metricAggregatorDTO: metricAggregatorDTO,
	}
}
