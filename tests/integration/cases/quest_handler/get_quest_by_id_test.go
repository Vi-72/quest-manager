package quest_handler

// HANDLER LAYER INTEGRATION TESTS
// Tests for getQuestByIDHandler.Handle orchestration logic

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

func (s *Suite) TestGetQuestByID() {
	ctx := context.Background()

	// Pre-condition - create quest
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - get quest by ID
	foundQuest, err := casesteps.GetQuestByIDStep(ctx, s.TestDIContainer.GetQuestByIDHandler, createdQuest.ID())

	// Assert
	s.Require().NoError(err)

	// Create assertions for single quest verification
	singleAssertions := assertions.NewQuestSingleAssertions(s.Assert())
	singleAssertions.QuestMatchesCreated(foundQuest, createdQuest)
	singleAssertions.QuestHasValidLocationData(foundQuest)
}

func (s *Suite) TestGetQuestByIDNotFound() {
	ctx := context.Background()

	// Pre-condition - use non-existent quest ID
	nonExistentID := uuid.New()

	// Act - try to get quest by non-existent ID
	_, err := casesteps.GetQuestByIDStep(ctx, s.TestDIContainer.GetQuestByIDHandler, nonExistentID)

	// Assert - should return error
	s.Require().Error(err, "Should return error for non-existent quest ID")
	s.Assert().Contains(err.Error(), "not found", "Error should indicate quest was not found")
}

func (s *Suite) TestGetQuestByIDHasAddresses() {
	ctx := context.Background()

	// Pre-condition - create quest with explicit different locations
	targetLocation := testdatagenerators.DefaultTestCoordinate()
	executionLocation := kernel.GeoCoordinate{
		Lat: targetLocation.Latitude() + 0.1,  // Much more significant difference
		Lon: targetLocation.Longitude() + 0.1, // Much more significant difference
	}

	// Get default test data and customize for addresses test
	questData := testdatagenerators.QuestDataWithLocations(targetLocation, executionLocation)

	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questData)
	s.Require().NoError(err)

	// Act - get quest by ID
	foundQuest, err := casesteps.GetQuestByIDStep(ctx, s.TestDIContainer.GetQuestByIDHandler, createdQuest.ID())

	// Assert
	s.Require().NoError(err)

	// Create assertions for location verification
	singleAssertions := assertions.NewQuestSingleAssertions(s.Assert())
	singleAssertions.QuestHasValidLocationData(foundQuest)
	s.Assert().Contains(*foundQuest.TargetAddress, "Target Address", "Target address should contain expected prefix")
	s.Assert().Contains(*foundQuest.ExecutionAddress, "Execution Address", "Execution address should contain expected prefix")
}
