package aggregator

import "time"

type Labels map[string]string

type MetricAggregatorDTO struct {
	ID                  string            `json:"id"`
	LifeSpan            int               `json:"life_span"`
	CreatedAt           time.Time         `json:"created_at"`
	Name                string            `json:"name"`
	Help                string            `json:"help"`
	ConstLabels         Labels            `json:"const_labels"`
	Buckets             []float64         `json:"buckets"`
	Stage               string            `json:"stage"`
	LabelValues         map[string]string `json:"label_values"`
	ForceRegistryUpdate bool              `json:"force_registry_update"`
	TimeToLiveSeconds   int               `json:"time_to_live"`
}
