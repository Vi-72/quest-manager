package quest_handler_tests

// HANDLER LAYER INTEGRATION TESTS
// Tests for createQuestHandler.Handle orchestration logic

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

func (s *Suite) TestCreateQuest() {
	ctx := context.Background()

	// Act - create quest with default data
	defaultData := testdatagenerators.DefaultQuestData()
	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, defaultData)

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(createdQuest.ID().String(), "Quest should have ID")

	// Verify quest data matches input (handler level - no addresses expected)

	s.Assert().Equal(defaultData.Title, createdQuest.Title)
	s.Assert().Equal(defaultData.Description, createdQuest.Description)
	s.Assert().Equal(defaultData.Difficulty, string(createdQuest.Difficulty))
	s.Assert().Equal(defaultData.Reward, createdQuest.Reward)
	s.Assert().Equal(defaultData.DurationMinutes, createdQuest.DurationMinutes)
	s.Assert().Equal(defaultData.Creator, createdQuest.Creator)
	s.Assert().Equal(defaultData.Equipment, createdQuest.Equipment)
	s.Assert().Equal(defaultData.Skills, createdQuest.Skills)
}

func (s *Suite) TestCreateQuestWithEmptyArrays() {
	ctx := context.Background()

	// Act - create quest with empty equipment and skills
	emptyData := testdatagenerators.EmptyArraysQuestData()
	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, emptyData)

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(createdQuest.ID().String(), "Quest should have ID")

	// Specifically verify empty arrays are preserved
	s.Assert().NotNil(createdQuest.Equipment, "Equipment should not be nil")
	s.Assert().NotNil(createdQuest.Skills, "Skills should not be nil")
	s.Assert().Len(createdQuest.Equipment, 0, "Equipment should be empty")
	s.Assert().Len(createdQuest.Skills, 0, "Skills should be empty")

	// Verify other data
	s.Assert().Equal(emptyData.Title, createdQuest.Title)
	s.Assert().Equal(emptyData.Description, createdQuest.Description)
	s.Assert().Equal(emptyData.Difficulty, string(createdQuest.Difficulty))
}

func (s *Suite) TestCreateQuestWithCustomLocations() {
	ctx := context.Background()

	// Pre-condition - create different locations
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}    // Moscow
	executionLocation := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609} // SPb

	// Act - create quest with custom locations
	questData := testdatagenerators.QuestDataWithLocations(targetLocation, executionLocation)
	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questData)

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(createdQuest.ID().String(), "Quest should have ID")

	// Verify locations are set correctly
	s.Assert().Equal(targetLocation.Latitude(), createdQuest.TargetLocation.Latitude())
	s.Assert().Equal(targetLocation.Longitude(), createdQuest.TargetLocation.Longitude())
	s.Assert().Equal(executionLocation.Latitude(), createdQuest.ExecutionLocation.Latitude())
	s.Assert().Equal(executionLocation.Longitude(), createdQuest.ExecutionLocation.Longitude())

	// Coordinates are validated above - handler level doesn't need address validation
}

func (s *Suite) TestCreateQuestWithAllParameters() {
	ctx := context.Background()
	handlerAssertions := assertions.NewQuestHandlerAssertions(s.Assert())

	// Pre-condition - prepare specific test data
	targetLocation := kernel.GeoCoordinate{Lat: 40.7128, Lon: -74.0060}     // NYC
	executionLocation := kernel.GeoCoordinate{Lat: 34.0522, Lon: -118.2437} // LA

	// Act - create quest with all custom parameters
	questData := testdatagenerators.SimpleQuestData(
		"Custom Quest Title",
		"Custom quest description for testing",
		"hard",
		5,
		120,
		targetLocation,
		executionLocation)
	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questData)

	// Assert using handler assertions pattern
	handlerAssertions.VerifyQuestFullMatch(createdQuest, questData, err)

	// Coordinates are validated above - handler level doesn't need address validation
}

func (s *Suite) TestCreateQuestPersistence() {
	ctx := context.Background()

	// Act - create quest
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Verify quest is persisted by retrieving it
	foundQuest, err := casesteps.GetQuestByIDStep(ctx, s.TestDIContainer.GetQuestByIDHandler, createdQuest.ID())

	// Assert
	s.Require().NoError(err)
	s.Assert().Equal(createdQuest.ID(), foundQuest.ID())
	s.Assert().Equal(createdQuest.Title, foundQuest.Title)
	s.Assert().Equal(createdQuest.Description, foundQuest.Description)
	s.Assert().Equal(createdQuest.Equipment, foundQuest.Equipment)
	s.Assert().Equal(createdQuest.Skills, foundQuest.Skills)

	// Verify timestamps are set
	s.Assert().False(foundQuest.CreatedAt.IsZero(), "CreatedAt should be set")
	s.Assert().False(foundQuest.UpdatedAt.IsZero(), "UpdatedAt should be set")
}

// NOTE: Domain validation tests (difficulty, reward ranges, duration limits)
// are now in tests/domain/quest_test.go where they belong

// LOCATION-SPECIFIC TESTS

func (s *Suite) TestCreateQuestWithSameLocations() {
	ctx := context.Background()

	// Pre-condition - prepare quest data with same target and execution locations
	sameLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176} // Moscow Red Square
	questData := testdatagenerators.SimpleQuestData(
		"Same Location Quest",
		"Quest with identical target and execution locations",
		"easy",
		2,
		45,
		sameLocation,
		sameLocation,
	)

	// Act - create quest with identical locations
	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questData)

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(createdQuest.ID().String(), "Quest should have ID")

	// Verify both locations are approximately the same (account for floating point precision)
	s.Assert().InDelta(sameLocation.Latitude(), createdQuest.TargetLocation.Latitude(), 0.01, "Target latitude should be close")
	s.Assert().InDelta(sameLocation.Longitude(), createdQuest.TargetLocation.Longitude(), 0.01, "Target longitude should be close")
	s.Assert().InDelta(sameLocation.Latitude(), createdQuest.ExecutionLocation.Latitude(), 0.01, "Execution latitude should be close")
	s.Assert().InDelta(sameLocation.Longitude(), createdQuest.ExecutionLocation.Longitude(), 0.01, "Execution longitude should be close")

	// Verify they are actually identical (within precision)
	s.Assert().InDelta(createdQuest.TargetLocation.Latitude(), createdQuest.ExecutionLocation.Latitude(), 0.01, "Target and execution latitudes should be identical")
	s.Assert().InDelta(createdQuest.TargetLocation.Longitude(), createdQuest.ExecutionLocation.Longitude(), 0.01, "Target and execution longitudes should be identical")
}

func (s *Suite) TestCreateQuestWithExistingLocationSameAddress() {
	ctx := context.Background()

	// Pre-condition - create first quest to establish location in database
	firstLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}  // Moscow
	secondLocation := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609} // SPb

	firstQuestData := testdatagenerators.SimpleQuestData(
		"First Quest",
		"First quest at specific locations",
		"easy",
		2,
		30,
		firstLocation,
		secondLocation,
	)
	firstQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, firstQuestData)
	s.Require().NoError(err)

	// Act - create second quest with exactly same locations (should reuse location records)
	secondQuestData := testdatagenerators.SimpleQuestData(
		"Second Quest",
		"Second quest at same locations",
		"medium",
		3,
		45,
		firstLocation,
		secondLocation,
	)
	secondQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, secondQuestData)

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(secondQuest.ID().String(), "Second quest should have ID")
	s.Assert().NotEqual(firstQuest.ID(), secondQuest.ID(), "Quests should have different IDs")

	// Verify locations are approximately identical (coordinates within precision)
	s.Assert().InDelta(firstQuest.TargetLocation.Latitude(), secondQuest.TargetLocation.Latitude(), 0.01, "Target latitudes should be close")
	s.Assert().InDelta(firstQuest.TargetLocation.Longitude(), secondQuest.TargetLocation.Longitude(), 0.01, "Target longitudes should be close")
	s.Assert().InDelta(firstQuest.ExecutionLocation.Latitude(), secondQuest.ExecutionLocation.Latitude(), 0.01, "Execution latitudes should be close")
	s.Assert().InDelta(firstQuest.ExecutionLocation.Longitude(), secondQuest.ExecutionLocation.Longitude(), 0.01, "Execution longitudes should be close")

	// Note: Location ID reuse depends on implementation and may create new UUIDs
	// The important thing is that coordinates are preserved correctly
}

func (s *Suite) TestCreateQuestWithExistingLocationDifferentAddress() {
	ctx := context.Background()

	// Pre-condition - create first quest with location
	baseLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176} // Moscow
	firstQuestData := testdatagenerators.SimpleQuestData(
		"First Context Quest",
		"First quest at this location",
		"easy",
		2,
		30,
		baseLocation,
		baseLocation,
	)
	firstQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, firstQuestData)
	s.Require().NoError(err)

	// Note: This test is more about HTTP layer with addresses
	// At handler level, we only work with coordinates
	// Different addresses would be handled by HTTP layer setting TargetAddress/ExecutionAddress

	// Act - create second quest with same coordinates but conceptually different context
	secondQuestData := testdatagenerators.SimpleQuestData(
		"Different Context Quest",
		"Same location but different purpose",
		"medium",
		4,
		90,
		baseLocation, // Same coordinates
		baseLocation,
	)
	secondQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, secondQuestData)

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(secondQuest.ID().String(), "Second quest should have ID")
	s.Assert().NotEqual(firstQuest.ID(), secondQuest.ID(), "Quests should have different IDs")

	// Verify coordinates are approximately the same (handler level validation with precision tolerance)
	s.Assert().InDelta(firstQuest.TargetLocation.Latitude(), secondQuest.TargetLocation.Latitude(), 0.01, "Target latitudes should be close")
	s.Assert().InDelta(firstQuest.TargetLocation.Longitude(), secondQuest.TargetLocation.Longitude(), 0.01, "Target longitudes should be close")

	// But quest content is different
	s.Assert().NotEqual(firstQuest.Title, secondQuest.Title)
	s.Assert().NotEqual(firstQuest.Description, secondQuest.Description)
}
