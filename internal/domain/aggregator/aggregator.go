package aggregator

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

var mutex = &sync.Mutex{}

type Aggregator interface {
	CollectHistogram(histogramAggregator MetricAggregatorDTO) (interface{}, error)
}

type aggregator struct {
	repository Repository
	collector  Collector
}

type Input struct {
	Repository Repository `validate:"required"`
	Collector  Collector  `validate:"required"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func New(input Input) (Aggregator, error) {

	err := validate.Struct(input)

	if err != nil {
		return nil, err
	}

	return &aggregator{
		repository: input.Repository,
		collector:  input.Collector,
	}, nil
}

func (a *aggregator) CollectHistogram(dto MetricAggregatorDTO) (interface{}, error) {
	return a.schedule(dto)
}

func (a *aggregator) schedule(dto MetricAggregatorDTO) (interface{}, error) {

	entity, err := a.findEntity(dto)
	if err != nil {
		return nil, err
	}

	go a.task(entity, a.collector)

	return "ok", nil
}

func (a *aggregator) task(entity *Entity, collector Collector) {
	duration := time.Since(entity.CreatedAt).Seconds()
	t0 := time.Now()
	lifeSpanDuration := time.Duration(entity.LastLifeSpan()) * time.Second
	time.Sleep(lifeSpanDuration)

	mutex.Lock()
	entity, err := a.repository.Find(entity.ID)
	mutex.Unlock()

	if err != nil {
		fmt.Printf("1>>>>>>>>>>>>>>>> err: %v\n", err)
		return
	}

	if !entity.InProgress() {
		return
	}

	if entity.UpdateAt.After(t0) {
		return
	}

	collector.Collect(entity, duration)

	mutex.Lock()
	entity.Close()
	_, err = a.repository.Update(*entity)
	mutex.Unlock()

	if err != nil {
		fmt.Printf("2>>>>>>>>>>>>>>>> err: %v\n", err)
		return
	}

	time.Sleep(time.Duration(entity.TimeToLiveSeconds) * time.Second)

	mutex.Lock()
	err = a.repository.Delete(entity.ID)
	mutex.Unlock()

	if err != nil {
		fmt.Printf("3>>>>>>>>>>>>>>>> err: %v\n", err)
		return
	}
}

func (a *aggregator) findEntity(dto MetricAggregatorDTO) (*Entity, error) {

	entity, err := a.repository.Find(dto.ID)

	if err != nil {
		if err.Error() != "id not found in database" {
			return nil, err
		}
		return a.saveNewEntity(dto)

	}
	return a.updateEntityIfStageIsNew(dto, entity)

}

func makeStage(dto MetricAggregatorDTO) Stage {
	return Stage{
		Stage:       dto.Stage,
		LabelValues: dto.LabelValues,
		CreatedAt:   dto.CreatedAt,
		LifeSpan:    dto.LifeSpan,
	}
}

func (a *aggregator) saveNewEntity(dto MetricAggregatorDTO) (*Entity, error) {
	var createdAt time.Time
	if dto.CreatedAt.IsZero() {
		createdAt = time.Now()
	} else {
		createdAt = dto.CreatedAt
	}

	var timeToLiveSeconds int

	if dto.TimeToLiveSeconds > 0 {
		timeToLiveSeconds = dto.TimeToLiveSeconds
	} else {
		timeToLiveSeconds = 30
	}

	return a.repository.Insert(Entity{
		ID:                dto.ID,
		Name:              dto.Name,
		Help:              dto.Help,
		CreatedAt:         createdAt,
		State:             "IN_PROGRESS",
		TimeToLiveSeconds: timeToLiveSeconds,
		Stages: []Stage{
			makeStage(dto),
		},
	})
}

func (a *aggregator) updateEntityIfStageIsNew(dto MetricAggregatorDTO, entity *Entity) (*Entity, error) {
	isNew := true
	for _, stage := range entity.Stages {
		if stage.Stage == dto.Stage {
			isNew = false
			break
		}
	}

	if !isNew {
		return entity, nil
	}

	entity.Stages = append(entity.Stages, makeStage(dto))
	entity.UpdateAt = time.Now()
	return a.repository.Update(*entity)
}
