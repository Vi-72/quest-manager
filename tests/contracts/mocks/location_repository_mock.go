package mocks

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"

	"github.com/google/uuid"
)

// MockLocationRepository is an in-memory implementation of LocationRepository for contract testing
type MockLocationRepository struct {
	locations map[uuid.UUID]*location.Location
	mu        sync.RWMutex
}

func NewMockLocationRepository() *MockLocationRepository {
	return &MockLocationRepository{
		locations: make(map[uuid.UUID]*location.Location),
	}
}

func (m *MockLocationRepository) GetByID(ctx context.Context, locationID uuid.UUID) (*location.Location, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	loc, exists := m.locations[locationID]
	if !exists {
		return nil, fmt.Errorf("location with id %s not found", locationID.String())
	}
	return loc, nil
}

func (m *MockLocationRepository) Save(ctx context.Context, l *location.Location) error {
	_ = ctx // unused in mock
	m.mu.Lock()
	defer m.mu.Unlock()
	m.locations[l.ID()] = l
	return nil
}

func (m *MockLocationRepository) FindAll(ctx context.Context) ([]*location.Location, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*location.Location
	for _, loc := range m.locations {
		result = append(result, loc)
	}
	return result, nil
}

func (m *MockLocationRepository) FindByBoundingBox(ctx context.Context, bbox kernel.BoundingBox) ([]*location.Location, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*location.Location
	for _, loc := range m.locations {
		if m.isWithinBoundingBox(loc.Coordinate, bbox) {
			result = append(result, loc)
		}
	}
	return result, nil
}

func (m *MockLocationRepository) FindByName(ctx context.Context, namePattern string) ([]*location.Location, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*location.Location
	namePattern = strings.ToLower(namePattern)

	for _, loc := range m.locations {
		if loc.Address != nil {
			if strings.Contains(strings.ToLower(*loc.Address), namePattern) {
				result = append(result, loc)
			}
		}
	}
	return result, nil
}

func (m *MockLocationRepository) isWithinBoundingBox(coord kernel.GeoCoordinate, bbox kernel.BoundingBox) bool {
	return coord.Lat >= bbox.MinLat && coord.Lat <= bbox.MaxLat &&
		coord.Lon >= bbox.MinLon && coord.Lon <= bbox.MaxLon
}

// Helper methods for testing
func (m *MockLocationRepository) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.locations = make(map[uuid.UUID]*location.Location)
}

func (m *MockLocationRepository) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.locations)
}
