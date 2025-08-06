package contracts

import (
	"context"
	"testing"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/tests/contracts/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// QuestRepositoryContractSuite defines contract tests that all QuestRepository implementations must pass
type QuestRepositoryContractSuite struct {
	suite.Suite
	repo ports.QuestRepository
	ctx  context.Context
}

// LocationRepositoryContractSuite defines contract tests that all LocationRepository implementations must pass
type LocationRepositoryContractSuite struct {
	suite.Suite
	repo ports.LocationRepository
	ctx  context.Context
}

func (s *QuestRepositoryContractSuite) SetupSuite() {
	s.repo = mocks.NewMockQuestRepository()
	s.ctx = context.Background()
}

func (s *QuestRepositoryContractSuite) SetupTest() {
	// Clear repository before each test
	if mockRepo, ok := s.repo.(*mocks.MockQuestRepository); ok {
		mockRepo.Clear()
	}
}

func (s *LocationRepositoryContractSuite) SetupSuite() {
	s.repo = mocks.NewMockLocationRepository()
	s.ctx = context.Background()
}

func (s *LocationRepositoryContractSuite) SetupTest() {
	// Clear repository before each test
	if mockRepo, ok := s.repo.(*mocks.MockLocationRepository); ok {
		mockRepo.Clear()
	}
}

// TestQuestRepositoryContract runs the contract test suite for QuestRepository
func TestQuestRepositoryContract(t *testing.T) {
	suite.Run(t, new(QuestRepositoryContractSuite))
}

// TestLocationRepositoryContract runs the contract test suite for LocationRepository
func TestLocationRepositoryContract(t *testing.T) {
	suite.Run(t, new(LocationRepositoryContractSuite))
}

// QuestRepository contract tests

func (s *QuestRepositoryContractSuite) TestSaveAndGetByID() {
	// Create a test quest
	coord1 := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}
	coord2 := kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0}

	q, err := quest.NewQuest(
		"Test Quest",
		"Description for contract test",
		"easy",
		2,
		30,
		coord1,
		coord2,
		"test-creator",
		[]string{},
		[]string{},
	)
	s.Require().NoError(err)

	// Contract: Save should store the quest without error
	err = s.repo.Save(s.ctx, q)
	s.Require().NoError(err, "Save operation should succeed")

	// Contract: GetByID should retrieve the saved quest
	retrieved, err := s.repo.GetByID(s.ctx, q.ID())
	s.Require().NoError(err, "GetByID operation should succeed")

	// Contract: Retrieved quest should match the original
	s.Assert().Equal(q.ID(), retrieved.ID(), "Quest ID should match")
	s.Assert().Equal(q.Title, retrieved.Title, "Quest title should match")
	s.Assert().Equal(q.Description, retrieved.Description, "Quest description should match")
	s.Assert().Equal(q.Difficulty, retrieved.Difficulty, "Quest difficulty should match")
}

func (s *QuestRepositoryContractSuite) TestGetByIDNonExistent() {
	// Contract: GetByID should return error for non-existent quest
	nonExistentID := uuid.New()

	_, err := s.repo.GetByID(s.ctx, nonExistentID)
	s.Assert().Error(err, "GetByID should return error for non-existent quest")
}

func (s *QuestRepositoryContractSuite) TestFindAll() {
	// Create multiple test quests
	coord1 := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}
	coord2 := kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0}

	q1, err := quest.NewQuest("Quest 1", "First quest", "easy", 2, 30, coord1, coord2, "creator1", []string{}, []string{})
	s.Require().NoError(err)

	q2, err := quest.NewQuest("Quest 2", "Second quest", "medium", 3, 45, coord1, coord2, "creator2", []string{}, []string{})
	s.Require().NoError(err)

	// Save both quests
	err = s.repo.Save(s.ctx, q1)
	s.Require().NoError(err)
	err = s.repo.Save(s.ctx, q2)
	s.Require().NoError(err)

	// Contract: FindAll should return all saved quests
	quests, err := s.repo.FindAll(s.ctx)
	s.Require().NoError(err, "FindAll operation should succeed")
	s.Assert().GreaterOrEqual(len(quests), 2, "Should return at least the saved quests")

	// Find our saved quests in the results
	found1, found2 := false, false
	for _, q := range quests {
		if q.ID() == q1.ID() {
			found1 = true
		}
		if q.ID() == q2.ID() {
			found2 = true
		}
	}
	s.Assert().True(found1, "Should find first saved quest")
	s.Assert().True(found2, "Should find second saved quest")
}

func (s *QuestRepositoryContractSuite) TestFindByStatus() {
	// Create a test quest with created status
	coord1 := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}
	coord2 := kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0}

	q, err := quest.NewQuest("Status Test Quest", "Quest for status testing", "easy", 2, 30, coord1, coord2, "creator", []string{}, []string{})
	s.Require().NoError(err)

	// Save the quest
	err = s.repo.Save(s.ctx, q)
	s.Require().NoError(err)

	// Contract: FindByStatus should return quests with the specified status
	quests, err := s.repo.FindByStatus(s.ctx, quest.StatusCreated)
	s.Require().NoError(err, "FindByStatus operation should succeed")

	// Find our quest in the results
	found := false
	for _, foundQuest := range quests {
		if foundQuest.ID() == q.ID() {
			found = true
			s.Assert().Equal(quest.StatusCreated, foundQuest.Status, "Quest should have created status")
			break
		}
	}
	s.Assert().True(found, "Should find quest with created status")
}

func (s *QuestRepositoryContractSuite) TestFindByBoundingBox() {
	// Create a test quest at specific coordinates
	coord := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}

	q, err := quest.NewQuest("Location Test Quest", "Quest for location testing", "easy", 2, 30, coord, coord, "creator", []string{}, []string{})
	s.Require().NoError(err)

	// Save the quest
	err = s.repo.Save(s.ctx, q)
	s.Require().NoError(err)

	// Contract: FindByBoundingBox should return quests within the bounding box
	bbox := kernel.BoundingBox{
		MinLat: 49.0,
		MaxLat: 51.0,
		MinLon: 9.0,
		MaxLon: 11.0,
	}

	quests, err := s.repo.FindByBoundingBox(s.ctx, bbox)
	s.Require().NoError(err, "FindByBoundingBox operation should succeed")

	// Find our quest in the results
	found := false
	for _, foundQuest := range quests {
		if foundQuest.ID() == q.ID() {
			found = true
			break
		}
	}
	s.Assert().True(found, "Should find quest within bounding box")
}

// LocationRepository contract tests

func (s *LocationRepositoryContractSuite) TestSaveAndGetByID() {
	// Create a test location
	coord := kernel.GeoCoordinate{Lat: 52.5200, Lon: 13.4050} // Berlin coordinates
	address := "Berlin, Germany"

	loc, err := location.NewLocation(coord, &address)
	s.Require().NoError(err)

	// Contract: Save should store the location without error
	err = s.repo.Save(s.ctx, loc)
	s.Require().NoError(err, "Save operation should succeed")

	// Contract: GetByID should retrieve the saved location
	retrieved, err := s.repo.GetByID(s.ctx, loc.ID())
	s.Require().NoError(err, "GetByID operation should succeed")

	// Contract: Retrieved location should match the original
	s.Assert().Equal(loc.ID(), retrieved.ID(), "Location ID should match")
	s.Assert().Equal(loc.Coordinate.Lat, retrieved.Coordinate.Lat, "Latitude should match")
	s.Assert().Equal(loc.Coordinate.Lon, retrieved.Coordinate.Lon, "Longitude should match")
	s.Assert().Equal(*loc.Address, *retrieved.Address, "Address should match")
}

func (s *LocationRepositoryContractSuite) TestGetByIDNonExistent() {
	// Contract: GetByID should return error for non-existent location
	nonExistentID := uuid.New()

	_, err := s.repo.GetByID(s.ctx, nonExistentID)
	s.Assert().Error(err, "GetByID should return error for non-existent location")
}

func (s *LocationRepositoryContractSuite) TestFindAll() {
	// Create multiple test locations
	coord1 := kernel.GeoCoordinate{Lat: 48.8566, Lon: 2.3522}  // Paris
	coord2 := kernel.GeoCoordinate{Lat: 51.5074, Lon: -0.1278} // London

	address1 := "Paris, France"
	address2 := "London, UK"

	loc1, err := location.NewLocation(coord1, &address1)
	s.Require().NoError(err)

	loc2, err := location.NewLocation(coord2, &address2)
	s.Require().NoError(err)

	// Save both locations
	err = s.repo.Save(s.ctx, loc1)
	s.Require().NoError(err)
	err = s.repo.Save(s.ctx, loc2)
	s.Require().NoError(err)

	// Contract: FindAll should return all saved locations
	locations, err := s.repo.FindAll(s.ctx)
	s.Require().NoError(err, "FindAll operation should succeed")
	s.Assert().GreaterOrEqual(len(locations), 2, "Should return at least the saved locations")

	// Find our saved locations in the results
	found1, found2 := false, false
	for _, l := range locations {
		if l.ID() == loc1.ID() {
			found1 = true
		}
		if l.ID() == loc2.ID() {
			found2 = true
		}
	}
	s.Assert().True(found1, "Should find first saved location")
	s.Assert().True(found2, "Should find second saved location")
}

func (s *LocationRepositoryContractSuite) TestFindByBoundingBox() {
	// Create a test location at specific coordinates
	coord := kernel.GeoCoordinate{Lat: 40.7128, Lon: -74.0060} // New York
	address := "New York, NY"

	loc, err := location.NewLocation(coord, &address)
	s.Require().NoError(err)

	// Save the location
	err = s.repo.Save(s.ctx, loc)
	s.Require().NoError(err)

	// Contract: FindByBoundingBox should return locations within the bounding box
	bbox := kernel.BoundingBox{
		MinLat: 40.0,
		MaxLat: 41.0,
		MinLon: -75.0,
		MaxLon: -73.0,
	}

	locations, err := s.repo.FindByBoundingBox(s.ctx, bbox)
	s.Require().NoError(err, "FindByBoundingBox operation should succeed")

	// Find our location in the results
	found := false
	for _, foundLoc := range locations {
		if foundLoc.ID() == loc.ID() {
			found = true
			s.Assert().Equal(coord.Lat, foundLoc.Coordinate.Lat, "Latitude should match")
			s.Assert().Equal(coord.Lon, foundLoc.Coordinate.Lon, "Longitude should match")
			break
		}
	}
	s.Assert().True(found, "Should find location within bounding box")
}
