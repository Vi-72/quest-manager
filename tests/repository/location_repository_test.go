package repository

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"
)

type LocationRepositoryTestSuite struct {
	RepositoryTestSuite
}

func (suite *LocationRepositoryTestSuite) SetupTest() {
	// Call parent setup to clean database
	suite.RepositoryTestSuite.SetupTest()
}

func (suite *LocationRepositoryTestSuite) TestSave_Success() {
	ctx := context.Background()

	// Arrange - create a valid location
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address, Moscow"

	loc, err := location.NewLocation(coordinate, &address)
	suite.Require().NoError(err)

	// Act - save location
	err = suite.locationRepo.Save(ctx, loc)

	// Assert
	suite.Require().NoError(err)

	// Verify location was saved by retrieving it
	saved, err := suite.locationRepo.GetByID(ctx, loc.ID())
	suite.Require().NoError(err)
	suite.assertLocationEquals(*loc, *saved)
}

func (suite *LocationRepositoryTestSuite) TestSave_Update() {
	ctx := context.Background()

	// Arrange - save initial location
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Original Address"

	loc, err := location.NewLocation(coordinate, &address)
	suite.Require().NoError(err)

	err = suite.locationRepo.Save(ctx, loc)
	suite.Require().NoError(err)

	// Modify location
	newCoordinate := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609}
	newAddress := "Updated Address"
	err = loc.Update(newCoordinate, &newAddress)
	suite.Require().NoError(err)

	// Act - save updated location
	err = suite.locationRepo.Save(ctx, loc)

	// Assert
	suite.Require().NoError(err)

	// Verify location was updated
	updated, err := suite.locationRepo.GetByID(ctx, loc.ID())
	suite.Require().NoError(err)
	suite.Equal(newCoordinate, updated.Coordinate)
	suite.NotNil(updated.Address)
	suite.Equal(newAddress, *updated.Address)
}

func (suite *LocationRepositoryTestSuite) TestSave_WithNilAddress() {
	ctx := context.Background()

	// Arrange - create location without address
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}

	loc, err := location.NewLocation(coordinate, nil)
	suite.Require().NoError(err)

	// Act - save location
	err = suite.locationRepo.Save(ctx, loc)

	// Assert
	suite.Require().NoError(err)

	// Verify location was saved
	saved, err := suite.locationRepo.GetByID(ctx, loc.ID())
	suite.Require().NoError(err)
	suite.Equal(coordinate, saved.Coordinate)
	suite.Nil(saved.Address)
}

func (suite *LocationRepositoryTestSuite) TestGetByID_Success() {
	ctx := context.Background()

	// Arrange - save a location
	coordinate := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	address := "Test Address"

	loc, err := location.NewLocation(coordinate, &address)
	suite.Require().NoError(err)

	err = suite.locationRepo.Save(ctx, loc)
	suite.Require().NoError(err)

	// Act - get location by ID
	found, err := suite.locationRepo.GetByID(ctx, loc.ID())

	// Assert
	suite.Require().NoError(err)
	suite.assertLocationEquals(*loc, *found)
}

func (suite *LocationRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()

	// Act - try to get non-existent location
	nonExistentID := uuid.New()
	_, err := suite.locationRepo.GetByID(ctx, nonExistentID)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "not found")
}

func (suite *LocationRepositoryTestSuite) TestFindAll_Success() {
	ctx := context.Background()

	// Arrange - save multiple locations
	coord1 := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	coord2 := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609}

	address1 := "Moscow Address"
	address2 := "St Petersburg Address"

	loc1, err := location.NewLocation(coord1, &address1)
	suite.Require().NoError(err)
	loc2, err := location.NewLocation(coord2, &address2)
	suite.Require().NoError(err)

	err = suite.locationRepo.Save(ctx, loc1)
	suite.Require().NoError(err)
	err = suite.locationRepo.Save(ctx, loc2)
	suite.Require().NoError(err)

	// Act - find all locations
	locations, err := suite.locationRepo.FindAll(ctx)

	// Assert
	suite.Require().NoError(err)
	suite.Len(locations, 2)

	// Verify both locations are present
	locationIDs := make(map[string]bool)
	for _, l := range locations {
		locationIDs[l.ID().String()] = true
	}
	suite.True(locationIDs[loc1.ID().String()])
	suite.True(locationIDs[loc2.ID().String()])
}

func (suite *LocationRepositoryTestSuite) TestFindAll_Empty() {
	ctx := context.Background()

	// Act - find all locations in empty database
	locations, err := suite.locationRepo.FindAll(ctx)

	// Assert
	suite.Require().NoError(err)
	suite.Len(locations, 0)
}

func (suite *LocationRepositoryTestSuite) TestFindByBoundingBox_Success() {
	ctx := context.Background()

	// Arrange - create locations at different coordinates
	moscowCenter := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	moscowSuburb := kernel.GeoCoordinate{Lat: 55.8000, Lon: 37.7000}
	stPetersburg := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609}

	address1 := "Moscow Center"
	address2 := "Moscow Suburb"
	address3 := "St Petersburg"

	loc1, err := location.NewLocation(moscowCenter, &address1)
	suite.Require().NoError(err)
	loc2, err := location.NewLocation(moscowSuburb, &address2)
	suite.Require().NoError(err)
	loc3, err := location.NewLocation(stPetersburg, &address3)
	suite.Require().NoError(err)

	err = suite.locationRepo.Save(ctx, loc1)
	suite.Require().NoError(err)
	err = suite.locationRepo.Save(ctx, loc2)
	suite.Require().NoError(err)
	err = suite.locationRepo.Save(ctx, loc3)
	suite.Require().NoError(err)

	// Act - find locations in Moscow area bounding box
	moscowBBox := kernel.BoundingBox{
		MinLat: 55.5, MaxLat: 56.0,
		MinLon: 37.0, MaxLon: 38.0,
	}
	moscowLocations, err := suite.locationRepo.FindByBoundingBox(ctx, moscowBBox)

	// Assert
	suite.Require().NoError(err)
	suite.Len(moscowLocations, 2) // Should find both Moscow locations, but not St Petersburg

	moscowLocationIDs := make(map[string]bool)
	for _, l := range moscowLocations {
		moscowLocationIDs[l.ID().String()] = true
	}
	suite.True(moscowLocationIDs[loc1.ID().String()])
	suite.True(moscowLocationIDs[loc2.ID().String()])
}

// Helper methods
func (suite *LocationRepositoryTestSuite) assertLocationEquals(expected, actual location.Location) {
	suite.Equal(expected.ID(), actual.ID())
	suite.Equal(expected.Coordinate, actual.Coordinate)

	if expected.Address == nil {
		suite.Nil(actual.Address)
	} else {
		suite.NotNil(actual.Address)
		suite.Equal(*expected.Address, *actual.Address)
	}

	// Note: CreatedAt and UpdatedAt might differ slightly due to timing
	// In a real test, you might want to assert they're "close enough"
	suite.False(actual.CreatedAt.IsZero(), "CreatedAt should be set")
	suite.False(actual.UpdatedAt.IsZero(), "UpdatedAt should be set")
}
