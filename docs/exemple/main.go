package main

import (
	"time"

	latencyaggregator "github.com/Espigah/latency-aggregator/pkg/prometheus/latency-aggregator"
	"github.com/google/uuid"
)

func main() {
	traceid := uuid.New().String()
	latencyaggregator.New().
		ID(traceid).
		URL("http://localhost:7070").
		LifeSpan(2).
		Name("order").
		Stage("service1").
		Help("order metric description").
		Push()

	time.Sleep(1 * time.Second)

	latencyaggregator.New().
		ID(traceid).
		URL("http://localhost:7070").
		LifeSpan(2).
		Name("order").
		Stage("service2").
		Push()

	time.Sleep(1 * time.Second)

	latencyaggregator.New().
		ID(traceid).
		URL("http://localhost:7070").
		LifeSpan(1).
		Name("order").
		Stage("service3").
		Push()
}
