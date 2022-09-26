package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HistogramVec struct {
	prometheus.HistogramVec
}

func NewPrometheusHistogramVec(name, description string, labels map[string]string, metricLabel []string) (c HistogramVec) {
	if len(labels) == 0 {
		labels = map[string]string{}
	}

	c.HistogramVec = *promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        name,
		Help:        description,
		ConstLabels: labels,
	}, metricLabel)
	return
}

func (c *HistogramVec) ObserveWith(labels map[string]string, v float64) {
	c.HistogramVec.With(labels).Observe(v)
}
