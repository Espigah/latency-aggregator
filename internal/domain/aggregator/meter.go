package aggregator

type Meter interface {
	Registry(duration float64)
}
