package domain

// LOCATION DOMAIN EVENTS UNIT TESTS
// Tests for location domain events: location.created and location.updated

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"
)

func TestLocation_NewLocation_DomainEvents(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address, Moscow"

	// Act - create new location
	loc, err := location.NewLocation(coordinate, &address)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, loc)

	// Check that domain event was raised
	events := loc.GetDomainEvents()
	assert.Len(t, events, 1, "NewLocation should raise one domain event (location.created)")

	// Verify event is of correct type (implementation-specific)
	// Note: Specific event verification depends on domain event implementation
}

func TestLocation_Update_DomainEvents(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address"

	// Create location
	loc, err := location.NewLocation(coordinate, &address)
	assert.NoError(t, err)

	// Clear creation events
	loc.ClearDomainEvents()
	assert.Len(t, loc.GetDomainEvents(), 0, "Events should be cleared")

	// Act - update location
	newCoordinate := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609}
	newAddress := "Updated Address"
	err = loc.Update(newCoordinate, &newAddress)

	// Assert
	assert.NoError(t, err)

	events := loc.GetDomainEvents()
	assert.Len(t, events, 1, "Update should raise one domain event (location.updated)")

	// Verify event is of correct type (implementation-specific)
	// Note: Specific event verification depends on domain event implementation
}

func TestLocation_ClearDomainEvents(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address"

	// Create location
	loc, err := location.NewLocation(coordinate, &address)
	assert.NoError(t, err)

	// Ensure there are some events
	events := loc.GetDomainEvents()
	assert.NotEmpty(t, events, "Location should have domain events after creation")

	// Act - clear events
	loc.ClearDomainEvents()

	// Assert - events should be cleared
	events = loc.GetDomainEvents()
	assert.Empty(t, events, "Events should be cleared after ClearDomainEvents()")
}

func TestLocation_MultipleUpdates_DomainEvents(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address"

	// Create location
	loc, err := location.NewLocation(coordinate, &address)
	assert.NoError(t, err)

	// Clear creation events
	loc.ClearDomainEvents()

	// Perform multiple updates
	for i := 0; i < 3; i++ {
		newCoordinate := kernel.GeoCoordinate{
			Lat: 55.7558 + float64(i)*0.001,
			Lon: 37.6176 + float64(i)*0.001,
		}
		newAddress := "Address " + string(rune('A'+i))

		err = loc.Update(newCoordinate, &newAddress)
		assert.NoError(t, err)

		// Each update should raise one event
		events := loc.GetDomainEvents()
		assert.Len(t, events, i+1, "Each update should add one domain event")
	}

	// Final count should be 3 events
	events := loc.GetDomainEvents()
	assert.Len(t, events, 3, "Should have 3 events after 3 updates")
}
