package latencyaggregator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Labels map[string]string

type HistogramAggregator interface {
	Push() error
	ID(string) HistogramAggregator
	URL(string) HistogramAggregator
	CreatedAt(time.Time) HistogramAggregator
	LifeSpan(int) HistogramAggregator
	Name(string) HistogramAggregator
	Help(string) HistogramAggregator
	LabelValues(map[string]string) HistogramAggregator
	TimeToLiveSeconds(int) HistogramAggregator
	Stage(string) HistogramAggregator
}

type histogramAggregatorBuilder struct {
	payload *histogramAggregator
}

type histogramAggregator struct {
	url                 string
	ID                  string            `json:"id"`
	LifeSpan            int               `json:"life_span"`
	CreatedAt           time.Time         `json:"created_at"`
	Name                string            `json:"name"`
	Help                string            `json:"help"`
	Buckets             []float64         `json:"buckets"`
	Stage               string            `json:"stage"`
	LabelValues         map[string]string `json:"label_values"`
	ForceRegistryUpdate bool              `json:"force_registry_update"`
	TimeToLiveSeconds   int               `json:"time_to_live"`
}

func New() HistogramAggregator {
	return &histogramAggregatorBuilder{
		payload: &histogramAggregator{},
	}

}

func (h *histogramAggregatorBuilder) URL(value string) HistogramAggregator {
	h.payload.url = value
	return h
}

func (h *histogramAggregatorBuilder) ID(value string) HistogramAggregator {
	h.payload.ID = value
	return h
}

func (h *histogramAggregatorBuilder) LifeSpan(value int) HistogramAggregator {
	h.payload.LifeSpan = value
	return h
}

func (h *histogramAggregatorBuilder) CreatedAt(value time.Time) HistogramAggregator {
	h.payload.CreatedAt = value
	return h
}

func (h *histogramAggregatorBuilder) Name(value string) HistogramAggregator {
	h.payload.Name = value
	return h
}

func (h *histogramAggregatorBuilder) Help(value string) HistogramAggregator {
	h.payload.Help = value
	return h
}

func (h *histogramAggregatorBuilder) LabelValues(value map[string]string) HistogramAggregator {
	h.payload.LabelValues = value
	return h
}

func (h *histogramAggregatorBuilder) TimeToLiveSeconds(value int) HistogramAggregator {
	h.payload.TimeToLiveSeconds = value
	return h
}

func (h *histogramAggregatorBuilder) Stage(value string) HistogramAggregator {
	h.payload.Stage = value
	return h
}

func (h *histogramAggregatorBuilder) Push() error {
	_, _, err := post(h.payload)
	return err
}

type PushResult struct{}

func post(ha *histogramAggregator) (*PushResult, *http.Response, error) {
	body, _ := json.Marshal(ha)
	payload := bytes.NewBuffer(body)

	result := new(PushResult)

	url := ha.url
	response, err := request(http.MethodPost, url, payload, &result)

	if err != nil {
		return nil, nil, err
	}
	return result, response, nil
}
