package aggregator

type Collector interface {
	Collect(*Entity, float64)
}
