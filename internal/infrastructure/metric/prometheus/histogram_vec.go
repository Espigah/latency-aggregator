package prometheus

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type histogramVecObject struct {
	histogram *prometheus.HistogramVec
	startedAt time.Time
}

type HistogramVec interface {
	Start() HistogramVec
	Finished(constLabels map[string]string)
	Success()
	Error(code string)
	With(prometheus.Labels) prometheus.Observer
	Observe(float64, prometheus.Labels)
	Unregister() bool
}

func newHistogramVec(namespace, name, description string, bukets []float64, constLabels map[string]string, labelNames []string) *prometheus.HistogramVec {

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   namespace,
		Name:        name,
		Help:        description,
		ConstLabels: constLabels,
		Buckets:     bukets,
	}, labelNames)

	return histogramVec
}

func NewHistogramVec(namespace, name, description string, bukets []float64, constLabels map[string]string, labelNames []string) (HistogramVec, error) {

	histogramCached, _ := metricStorage.Load(name)

	if histogramCached != nil {
		return histogramCached.(HistogramVec), nil
	}

	histogramVec := newHistogramVec(namespace, name, description, bukets, constLabels, labelNames)

	err := prometheus.Register(histogramVec)

	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			fmt.Printf("are: %v\n", are)
			return &histogramVecObject{
				histogram: histogramVec,
			}, err
		}
		return nil, err
	}

	histogra := &histogramVecObject{
		histogram: histogramVec,
	}

	metricStorage.Store(name, histogra)

	return histogra, nil
}

func (h *histogramVecObject) Unregister() bool {
	return prometheus.Unregister(h.histogram)
}

func (h *histogramVecObject) Start() HistogramVec {
	return &histogramVecObject{
		startedAt: time.Now(),
		histogram: h.histogram,
	}
}

func (h *histogramVecObject) Finished(labels map[string]string) {
	duration := time.Since(h.startedAt).Seconds()
	h.Observe(duration, labels)
}

func (h *histogramVecObject) Observe(duration float64, labels prometheus.Labels) {
	h.With(labels).Observe(duration)
}

func (h *histogramVecObject) With(labels prometheus.Labels) prometheus.Observer {
	return h.histogram.With(labels)
}

func (h *histogramVecObject) Success() {
	h.Finished(prometheus.Labels{"status": "success", "status_code": "200"})
}

func (h *histogramVecObject) Error(code string) {
	h.Finished(prometheus.Labels{"status": "error", "status_code": code})
}
