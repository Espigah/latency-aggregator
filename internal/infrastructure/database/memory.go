package database

import (
	"errors"

	"github.com/Espigah/latency-aggregator/internal/domain/aggregator"
)

// NewMemoryDatabase create new memory database
func NewMemoryDatabase() aggregator.Repository {
	return &memoryDatabase{
		records: make(map[string]*aggregator.Entity),
	}
}

type memoryDatabase struct {
	records map[string]*aggregator.Entity
}

func (m *memoryDatabase) Find(id string) (*aggregator.Entity, error) {
	record := m.records[id]

	if record == nil {
		return nil, errors.New("id not found in database")
	}

	return record, nil
}

func (m *memoryDatabase) Insert(aggregatorEntity aggregator.Entity) (*aggregator.Entity, error) {
	m.records[aggregatorEntity.ID] = &aggregatorEntity

	return &aggregatorEntity, nil
}

func (m *memoryDatabase) Update(aggregatorEntity aggregator.Entity) (*aggregator.Entity, error) {
	return m.Insert(aggregatorEntity)
}

func (m *memoryDatabase) Delete(id string) error {
	delete(m.records, id)

	return nil
}
