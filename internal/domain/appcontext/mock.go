package appcontext

import (
	"context"

	"github.com/Espigah/latency-aggregator/internal/infrastructure/logger/logwrapper"
)

type mock struct {
	logger logwrapper.LoggerWrapper
}

// Mock returns a mock app context
func Mock(logger logwrapper.LoggerWrapper) Context {
	m := &mock{}
	m.SetLogger(logger)
	return m
}

func (mock *mock) SetLogger(logger logwrapper.LoggerWrapper) {
	mock.logger = logger
}

func (mock *mock) Logger() logwrapper.LoggerWrapper {
	return mock.logger
}

func (mock *mock) Context() context.Context {
	return context.Background()
}

func (mock *mock) TraceID() string {
	return ""
}

func (mock *mock) SpanID() string {
	return ""
}

func (mock *mock) Done() {
}

func (mock *mock) WithValue(key, val interface{}) {

}

func (mock *mock) Value(key interface{}) interface{} {
	return ""
}

func (mock *mock) TTL() *int64 {
	return nil
}
