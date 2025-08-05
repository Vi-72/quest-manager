package domain

// DOMAIN LAYER UNIT TESTS
// Tests for domain model business rules and validation logic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"
)

func TestNewLocation_Success(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address, Moscow"

	loc, err := location.NewLocation(coordinate, &address)

	assert.NoError(t, err)
	assert.NotNil(t, loc)
	assert.Equal(t, coordinate, loc.Coordinate)
	assert.NotNil(t, loc.Address)
	assert.Equal(t, address, *loc.Address)
	assert.False(t, loc.CreatedAt.IsZero())
	assert.False(t, loc.UpdatedAt.IsZero())
	assert.Equal(t, loc.CreatedAt, loc.UpdatedAt) // Should be same on creation
	assert.NotNil(t, loc.ID())

	// Check that domain event was raised
	events := loc.GetDomainEvents()
	assert.Len(t, events, 1, "NewLocation should raise one domain event")
}

func TestNewLocation_WithNilAddress(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}

	loc, err := location.NewLocation(coordinate, nil)

	assert.NoError(t, err)
	assert.NotNil(t, loc)
	assert.Equal(t, coordinate, loc.Coordinate)
	assert.Nil(t, loc.Address)
	assert.False(t, loc.CreatedAt.IsZero())
	assert.False(t, loc.UpdatedAt.IsZero())
	assert.NotNil(t, loc.ID())
}

func TestNewLocation_DifferentCoordinates(t *testing.T) {
	testCases := []struct {
		name       string
		coordinate kernel.GeoCoordinate
		address    string
	}{
		{
			name:       "Moscow center",
			coordinate: kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176},
			address:    "Red Square, Moscow",
		},
		{
			name:       "New York",
			coordinate: kernel.GeoCoordinate{Lat: 40.7128, Lon: -74.0060},
			address:    "Times Square, New York",
		},
		{
			name:       "North Pole",
			coordinate: kernel.GeoCoordinate{Lat: 90.0, Lon: 0.0},
			address:    "North Pole",
		},
		{
			name:       "Equator",
			coordinate: kernel.GeoCoordinate{Lat: 0.0, Lon: 0.0},
			address:    "Null Island",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loc, err := location.NewLocation(tc.coordinate, &tc.address)

			assert.NoError(t, err)
			assert.Equal(t, tc.coordinate, loc.Coordinate)
			assert.Equal(t, tc.address, *loc.Address)
		})
	}
}

func TestLocation_Update_Success(t *testing.T) {
	// Create initial location
	originalCoordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	originalAddress := "Original Address"

	loc, err := location.NewLocation(originalCoordinate, &originalAddress)
	assert.NoError(t, err)

	originalUpdatedAt := loc.UpdatedAt
	originalID := loc.ID()

	// Clear domain events from creation
	loc.ClearDomainEvents()

	// Small delay to ensure UpdatedAt changes
	time.Sleep(1 * time.Millisecond)

	// Update location
	newCoordinate := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609} // St. Petersburg
	newAddress := "New Address, St. Petersburg"

	err = loc.Update(newCoordinate, &newAddress)

	assert.NoError(t, err)
	assert.Equal(t, newCoordinate, loc.Coordinate)
	assert.Equal(t, newAddress, *loc.Address)
	assert.True(t, loc.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be updated")
	assert.Equal(t, originalID, loc.ID(), "ID should not change")

	// Check that domain event was raised
	events := loc.GetDomainEvents()
	assert.Len(t, events, 1, "Update should raise one domain event")
}

func TestLocation_Update_WithNilAddress(t *testing.T) {
	// Create location with address
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Original Address"

	loc, err := location.NewLocation(coordinate, &address)
	assert.NoError(t, err)

	// Update to nil address
	newCoordinate := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609}

	err = loc.Update(newCoordinate, nil)

	assert.NoError(t, err)
	assert.Equal(t, newCoordinate, loc.Coordinate)
	assert.Nil(t, loc.Address)
}

func TestLocation_Update_OnlyCoordinate(t *testing.T) {
	// Create location
	originalCoordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Unchanged Address"

	loc, err := location.NewLocation(originalCoordinate, &address)
	assert.NoError(t, err)

	// Update only coordinate, keep same address
	newCoordinate := kernel.GeoCoordinate{Lat: 55.7600, Lon: 37.6200}

	err = loc.Update(newCoordinate, &address)

	assert.NoError(t, err)
	assert.Equal(t, newCoordinate, loc.Coordinate)
	assert.Equal(t, address, *loc.Address)
}

func TestLocation_Update_OnlyAddress(t *testing.T) {
	// Create location
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	originalAddress := "Original Address"

	loc, err := location.NewLocation(coordinate, &originalAddress)
	assert.NoError(t, err)

	// Update only address, keep same coordinate
	newAddress := "Updated Address"

	err = loc.Update(coordinate, &newAddress)

	assert.NoError(t, err)
	assert.Equal(t, coordinate, loc.Coordinate)
	assert.Equal(t, newAddress, *loc.Address)
}

func TestLocation_DomainEvents(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address"

	// Test creation event
	loc, err := location.NewLocation(coordinate, &address)
	assert.NoError(t, err)

	events := loc.GetDomainEvents()
	assert.Len(t, events, 1, "NewLocation should raise LocationCreated event")

	// Clear events
	loc.ClearDomainEvents()
	assert.Len(t, loc.GetDomainEvents(), 0, "Events should be cleared")

	// Test update event
	newCoordinate := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609}
	newAddress := "Updated Address"

	err = loc.Update(newCoordinate, &newAddress)
	assert.NoError(t, err)

	events = loc.GetDomainEvents()
	assert.Len(t, events, 1, "Update should raise LocationUpdated event")
}

func TestLocation_ImmutableID(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address"

	loc, err := location.NewLocation(coordinate, &address)
	assert.NoError(t, err)

	originalID := loc.ID()

	// Update multiple times
	for i := 0; i < 5; i++ {
		newCoordinate := kernel.GeoCoordinate{
			Lat: 55.7558 + float64(i)*0.001,
			Lon: 37.6176 + float64(i)*0.001,
		}
		newAddress := "Address " + string(rune('A'+i))

		err = loc.Update(newCoordinate, &newAddress)
		assert.NoError(t, err)
		assert.Equal(t, originalID, loc.ID(), "ID should remain the same after update %d", i)
	}
}

func TestLocation_TimestampBehavior(t *testing.T) {
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address"

	loc, err := location.NewLocation(coordinate, &address)
	assert.NoError(t, err)

	createdAt := loc.CreatedAt
	originalUpdatedAt := loc.UpdatedAt

	// CreatedAt should equal UpdatedAt on creation
	assert.Equal(t, createdAt, originalUpdatedAt, "CreatedAt should equal UpdatedAt on creation")

	// Small delay to ensure time difference
	time.Sleep(2 * time.Millisecond)

	// Update location
	newCoordinate := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609}
	err = loc.Update(newCoordinate, &address)
	assert.NoError(t, err)

	// CreatedAt should not change, UpdatedAt should change
	assert.Equal(t, createdAt, loc.CreatedAt, "CreatedAt should not change after update")
	assert.True(t, loc.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should change after update")
	assert.True(t, loc.UpdatedAt.After(createdAt), "UpdatedAt should be after CreatedAt")
}
