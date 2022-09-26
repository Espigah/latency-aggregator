package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type gougeVecObject struct {
	gouge *prometheus.GaugeVec
}

type GaugeVec interface {
	Observe(float64, prometheus.Labels)
	Unregister() bool
}

func newGaugeVec(namespace, name, description string, bukets []float64, constLabels map[string]string, labelNames []string) *prometheus.GaugeVec {

	gougeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        name,
		Help:        description,
		ConstLabels: constLabels,
	}, labelNames)

	return gougeVec
}

func NewGaugeVec(namespace, name, description string, bukets []float64, constLabels map[string]string, labelNames []string) (GaugeVec, error) {

	gougeCached, _ := metricStorage.Load(name)

	if gougeCached != nil {
		return gougeCached.(GaugeVec), nil
	}

	gougeVec := newGaugeVec(namespace, name, description, bukets, constLabels, labelNames)

	err := prometheus.Register(gougeVec)

	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			fmt.Printf("are: %v\n", are)
			return &gougeVecObject{
				gouge: gougeVec,
			}, err
		}
		return nil, err
	}

	histogra := &gougeVecObject{
		gouge: gougeVec,
	}

	metricStorage.Store(name, histogra)

	return histogra, nil
}

func (h *gougeVecObject) Unregister() bool {
	return prometheus.Unregister(h.gouge)
}

func (h *gougeVecObject) Observe(duration float64, labels prometheus.Labels) {
	h.gouge.With(labels).Set(duration)
}
