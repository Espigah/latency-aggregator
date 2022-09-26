package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type GaugeVec struct {
	prometheus.GaugeVec
}

func NewPrometheusGaugeVec(name, description string, labels map[string]string, metricLabel []string) (c GaugeVec) {
	if len(labels) == 0 {
		labels = map[string]string{}
	}

	c.GaugeVec = *promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        name,
		Help:        description,
		ConstLabels: labels,
	}, metricLabel)
	return
}

func (c *GaugeVec) AddWith(labels map[string]string, v float64) {
	c.GaugeVec.With(labels).Add(v)
}
