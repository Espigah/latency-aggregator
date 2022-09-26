package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Histogram struct {
	prometheus.Histogram
}

func NewHistogram(name, description string, labels map[string]string) (c Histogram, err error) {
	if len(labels) == 0 {
		labels = map[string]string{}
	}

	c.Histogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:        name,
		Help:        description,
		ConstLabels: labels,
		Buckets:     []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 2.0, 3.0},
	})
	return
}

func (c *Histogram) Observe(v float64) {
	c.Histogram.Observe(v)
}
