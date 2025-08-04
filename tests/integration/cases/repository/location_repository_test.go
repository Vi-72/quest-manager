//go:build integration

package repository

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"

	"github.com/google/uuid"
)

func (s *Suite) TestLocationRepository_Save_Success() {
	ctx := context.Background()

	// Pre-condition - create a valid location
	loc := s.createTestLocation("Test Location", 55.7558, 37.6173)

	// Act - save location
	err := s.TestDIContainer.LocationRepository.Save(ctx, loc)

	// Assert
	s.Require().NoError(err)

	// Verify location was saved by retrieving it
	saved, err := s.TestDIContainer.LocationRepository.GetByID(ctx, loc.ID())
	s.Require().NoError(err)
	s.assertLocationEquals(*loc, *saved)
}

func (s *Suite) TestLocationRepository_Save_Update() {
	ctx := context.Background()

	// Pre-condition - save initial location
	loc := s.createTestLocation("Original Address", 55.7558, 37.6173)
	err := s.TestDIContainer.LocationRepository.Save(ctx, loc)
	s.Require().NoError(err)

	// Modify location
	newAddress := "Updated Address"
	newCoordinate := kernel.GeoCoordinate{Lat: 55.7600, Lon: 37.6200}

	err = loc.Update(newCoordinate, &newAddress)
	s.Require().NoError(err)

	// Act - save updated location
	err = s.TestDIContainer.LocationRepository.Save(ctx, loc)

	// Assert
	s.Require().NoError(err)

	// Verify location was updated
	updated, err := s.TestDIContainer.LocationRepository.GetByID(ctx, loc.ID())
	s.Require().NoError(err)
	s.Equal(newAddress, *updated.Address)
	s.Equal(newCoordinate, updated.Coordinate)
}

func (s *Suite) TestLocationRepository_Save_WithNilAddress() {
	ctx := context.Background()

	// Pre-condition - create location without address
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6173}
	loc, err := location.NewLocation(coordinate, nil)
	s.Require().NoError(err)

	// Act - save location
	err = s.TestDIContainer.LocationRepository.Save(ctx, loc)

	// Assert
	s.Require().NoError(err)

	// Verify location was saved
	saved, err := s.TestDIContainer.LocationRepository.GetByID(ctx, loc.ID())
	s.Require().NoError(err)
	s.Nil(saved.Address)
	s.Equal(coordinate, saved.Coordinate)
}

func (s *Suite) TestLocationRepository_GetByID_Success() {
	ctx := context.Background()

	// Pre-condition - save a location
	loc := s.createTestLocation("Test Location", 59.9311, 30.3609)
	err := s.TestDIContainer.LocationRepository.Save(ctx, loc)
	s.Require().NoError(err)

	// Act - get location by ID
	found, err := s.TestDIContainer.LocationRepository.GetByID(ctx, loc.ID())

	// Assert
	s.Require().NoError(err)
	s.assertLocationEquals(*loc, *found)
}

func (s *Suite) TestLocationRepository_GetByID_NotFound() {
	ctx := context.Background()

	// Pre-condition - use non-existent ID
	nonExistentID := uuid.New()

	// Act - try to get location by non-existent ID
	_, err := s.TestDIContainer.LocationRepository.GetByID(ctx, nonExistentID)

	// Assert - should return error
	s.Require().Error(err)
	s.Contains(err.Error(), "not found")
}

func (s *Suite) TestLocationRepository_FindAll_Empty() {
	ctx := context.Background()

	// Act - find all locations when database is empty
	locations, err := s.TestDIContainer.LocationRepository.FindAll(ctx)

	// Assert
	s.Require().NoError(err)
	s.Empty(locations)
}

func (s *Suite) TestLocationRepository_FindAll_Success() {
	ctx := context.Background()

	// Pre-condition - save multiple locations
	loc1 := s.createTestLocation("Location 1", 55.7558, 37.6173) // Moscow
	loc2 := s.createTestLocation("Location 2", 59.9311, 30.3609) // SPB
	loc3 := s.createTestLocation("Location 3", 51.5074, -0.1278) // London

	err := s.TestDIContainer.LocationRepository.Save(ctx, loc1)
	s.Require().NoError(err)
	err = s.TestDIContainer.LocationRepository.Save(ctx, loc2)
	s.Require().NoError(err)
	err = s.TestDIContainer.LocationRepository.Save(ctx, loc3)
	s.Require().NoError(err)

	// Act - find all locations
	locations, err := s.TestDIContainer.LocationRepository.FindAll(ctx)

	// Assert
	s.Require().NoError(err)
	s.Len(locations, 3)

	// Verify all locations are present
	locationIDs := make(map[uuid.UUID]bool)
	for _, l := range locations {
		locationIDs[l.ID()] = true
	}
	s.True(locationIDs[loc1.ID()])
	s.True(locationIDs[loc2.ID()])
	s.True(locationIDs[loc3.ID()])
}

func (s *Suite) TestLocationRepository_FindByBoundingBox_Success() {
	ctx := context.Background()

	// Pre-condition - create locations at different coordinates
	moscowLoc := s.createTestLocation("Moscow", 55.7558, 37.6173)
	spbLoc := s.createTestLocation("Saint Petersburg", 59.9311, 30.3609)
	londonLoc := s.createTestLocation("London", 51.5074, -0.1278)

	err := s.TestDIContainer.LocationRepository.Save(ctx, moscowLoc)
	s.Require().NoError(err)
	err = s.TestDIContainer.LocationRepository.Save(ctx, spbLoc)
	s.Require().NoError(err)
	err = s.TestDIContainer.LocationRepository.Save(ctx, londonLoc)
	s.Require().NoError(err)

	// Act - find locations in Russia area (bounding box)
	russiaBoundingBox := kernel.BoundingBox{
		MinLat: 55.0000, // South boundary (covers Moscow and SPB)
		MaxLat: 60.0000, // North boundary
		MinLon: 30.0000, // West boundary
		MaxLon: 40.0000, // East boundary
	}
	russianLocations, err := s.TestDIContainer.LocationRepository.FindByBoundingBox(ctx, russiaBoundingBox)

	// Assert
	s.Require().NoError(err)
	s.Len(russianLocations, 2) // Should find Moscow and SPB

	foundIDs := make(map[uuid.UUID]bool)
	for _, l := range russianLocations {
		foundIDs[l.ID()] = true
	}
	s.True(foundIDs[moscowLoc.ID()])
	s.True(foundIDs[spbLoc.ID()])
	s.False(foundIDs[londonLoc.ID()]) // London should not be found
}

// ==========================================
// POSTGRESQL-SPECIFIC TESTS
// ==========================================

func (s *Suite) TestLocationRepository_PostgreSQL_CoordinatePrecision() {
	ctx := context.Background()

	// Test PostgreSQL decimal precision for coordinates
	preciseCoordinate := kernel.GeoCoordinate{
		Lat: 55.758392847392847, // High precision latitude
		Lon: 37.617384738473847, // High precision longitude
	}

	address := "High Precision Location"
	loc, err := location.NewLocation(preciseCoordinate, &address)
	s.Require().NoError(err)

	// Save and retrieve
	err = s.TestDIContainer.LocationRepository.Save(ctx, loc)
	s.Require().NoError(err)

	found, err := s.TestDIContainer.LocationRepository.GetByID(ctx, loc.ID())
	s.Require().NoError(err)

	// Verify high precision is preserved in PostgreSQL
	s.InDelta(preciseCoordinate.Lat, found.Coordinate.Lat, 0.000001)
	s.InDelta(preciseCoordinate.Lon, found.Coordinate.Lon, 0.000001)
}

func (s *Suite) TestLocationRepository_PostgreSQL_UnicodeAddresses() {
	ctx := context.Background()

	// Test PostgreSQL handling of Unicode characters
	unicodeAddresses := []string{
		"Москва, Красная площадь, 1", // Russian
		"北京市天安门广场",                   // Chinese
		"東京都千代田区",                    // Japanese
		"🏠 Emoji Address 🌍",          // Emojis
	}

	savedLocations := make([]*location.Location, len(unicodeAddresses))

	for i, addr := range unicodeAddresses {
		coordinate := kernel.GeoCoordinate{Lat: 55.0 + float64(i), Lon: 37.0 + float64(i)}
		loc, err := location.NewLocation(coordinate, &addr)
		s.Require().NoError(err)

		err = s.TestDIContainer.LocationRepository.Save(ctx, loc)
		s.Require().NoError(err)

		savedLocations[i] = loc
	}

	// Verify all Unicode addresses are preserved
	for i, loc := range savedLocations {
		found, err := s.TestDIContainer.LocationRepository.GetByID(ctx, loc.ID())
		s.Require().NoError(err)
		s.Equal(unicodeAddresses[i], *found.Address)
	}
}

func (s *Suite) TestLocationRepository_PostgreSQL_Transactions() {
	ctx := context.Background()

	// Test that location changes within transaction are isolated
	loc1 := s.createTestLocation("Transaction Test 1", 55.7558, 37.6173)
	loc2 := s.createTestLocation("Transaction Test 2", 59.9311, 30.3609)

	// Save both locations
	err := s.TestDIContainer.LocationRepository.Save(ctx, loc1)
	s.Require().NoError(err)
	err = s.TestDIContainer.LocationRepository.Save(ctx, loc2)
	s.Require().NoError(err)

	// Both should be findable within this transaction
	locations, err := s.TestDIContainer.LocationRepository.FindAll(ctx)
	s.Require().NoError(err)
	s.Len(locations, 2)

	// Transaction will be rolled back in TearDownTest
	// so in next test they won't be there (tested implicitly)
}

// ==========================================
// HELPER METHODS
// ==========================================

func (s *Suite) createTestLocation(address string, lat, lon float64) *location.Location {
	coordinate := kernel.GeoCoordinate{Lat: lat, Lon: lon}
	loc, err := location.NewLocation(coordinate, &address)
	s.Require().NoError(err)
	return loc
}

func (s *Suite) assertLocationEquals(expected, actual location.Location) {
	s.Equal(expected.ID(), actual.ID())
	s.Equal(expected.Coordinate, actual.Coordinate)

	if expected.Address == nil {
		s.Nil(actual.Address)
	} else {
		s.Require().NotNil(actual.Address)
		s.Equal(*expected.Address, *actual.Address)
	}
}
