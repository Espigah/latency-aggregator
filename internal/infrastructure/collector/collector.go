package collector

import (
	"github.com/Espigah/latency-aggregator/internal/domain/aggregator"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/metric"
	"github.com/Espigah/latency-aggregator/internal/infrastructure/metric/prometheus"
)

type collector struct {
	aggregator.Repository
}

func New() aggregator.Collector {
	return &collector{}
}

func (c *collector) Collect(entity *aggregator.Entity, duration float64) {

	histogramCollector, err := c.NewHistogram(entity)

	if err != nil {
		return
	}

	histogramCollector.Registry(duration)

	gougeCollector, err := c.NewGouge(entity)

	if err != nil {
		return
	}

	gougeCollector.Registry(duration)

}

func (c *collector) NewGouge(entity *aggregator.Entity) (aggregator.Meter, error) {

	labelValues := entity.LabelValues()
	labelNames := makeLabelNames(labelValues)

	gougeVec, err := prometheus.NewBuilder().
		Name(entity.Name + "_total").
		Help(entity.Help + "(gouge)").
		LabelNames(labelNames).
		BuildGougeVec()

	if err != nil {
		return nil, err
	}

	gougeMetric := metric.NewGouge(
		gougeVec,
		aggregator.MetricAggregatorDTO{
			LabelValues: labelValues,
		},
	)

	return gougeMetric, nil
}

func (c *collector) NewHistogram(entity *aggregator.Entity) (aggregator.Meter, error) {

	labelValues := entity.LabelValues()
	labelNames := makeLabelNames(labelValues)

	histogramVec, err := prometheus.NewBuilder().
		Name(entity.Name).
		Help(entity.Help).
		LabelNames(labelNames).
		BuildHistogramVec()

	if err != nil {
		return nil, err
	}

	histogramMetric := metric.NewHistogram(
		histogramVec,
		aggregator.MetricAggregatorDTO{
			LabelValues: labelValues,
		},
	)

	return histogramMetric, nil
}

func makeLabelNames(labelValues map[string]string) []string {
	labelNames := make([]string, len(labelValues))

	index := 0
	for key := range labelValues {
		labelNames[index] = key
		index++
	}

	return labelNames
}
