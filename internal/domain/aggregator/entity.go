package aggregator

import "time"

type Entity struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Help              string    `json:"help"`
	CreatedAt         time.Time `json:"created_at"`
	UpdateAt          time.Time `json:"update_at"`
	Stages            []Stage   `json:"stages"`
	State             string    `json:"state"`
	TimeToLiveSeconds int       `json:"time_to_live"`
}

func (e *Entity) LabelValues() map[string]string {
	labels := make(map[string]string)
	for _, stage := range e.Stages {
		for key, value := range stage.LabelValues {
			labels[key] = value
		}
	}
	return labels
}

func (e *Entity) LastLifeSpan() int {
	return e.Stages[len(e.Stages)-1].LifeSpan
}

func (e *Entity) AddStage(stage Stage) {
	e.Stages = append(e.Stages, stage)
}

type Stage struct {
	Stage       string    `json:"stage"`
	CreatedAt   time.Time `json:"created_at"`
	LifeSpan    int       `json:"life_span"`
	LabelValues map[string]string
}

func (e *Stage) Equal(target Stage) bool {
	return e.Stage == target.Stage
}

func (e *Entity) InProgress() bool {
	return e.State != "closed"
}

func (e *Entity) Close() {
	e.State = "closed"
	e.UpdateAt = time.Now()
}
