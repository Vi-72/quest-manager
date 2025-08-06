package quest_http_tests

// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateUUID function

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

func (s *Suite) TestGetQuestByIDHTTP() {
	ctx := context.Background()

	// Pre-condition - create quest via handler (for setup)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - get quest by ID via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.ID().String())
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	foundQuest := httpAssertions.QuestHTTPGetSuccessfully(getResp, err)

	// Verify quest data matches created quest
	singleAssertions := assertions.NewQuestSingleAssertions(s.Assert())
	singleAssertions.QuestHTTPMatchesDomain(foundQuest, createdQuest)
	singleAssertions.QuestHTTPHasValidLocationData(foundQuest)
}

func (s *Suite) TestGetQuestByIDHTTPNotFound() {
	ctx := context.Background()

	// Pre-condition - use non-existent quest ID
	nonExistentID := uuid.New().String()

	// Act - try to get quest by non-existent ID via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(nonExistentID)
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert - should return 404 error
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	httpAssertions.QuestHTTPErrorResponse(getResp, err, http.StatusNotFound, "not found")
}

func (s *Suite) TestGetQuestByIDHTTPInvalidID() {
	ctx := context.Background()

	// Act - try to get quest with invalid UUID format via HTTP API
	getReq := casesteps.GetQuestHTTPRequest("invalid-uuid-format")
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert - should return 400 error for invalid UUID
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, getResp.StatusCode, "Should return 400 for invalid UUID format")

	// Verify error message contains validation details
	s.Assert().Contains(getResp.Body, "validation failed", "Error message should contain validation failure details")
	s.Assert().Contains(getResp.Body, "UUID", "Error message should mention UUID format requirement")
}

func (s *Suite) TestGetQuestByIDHTTPHasAddresses() {
	ctx := context.Background()

	// Pre-condition - create quest with explicit different locations via handler
	targetLocation := testdatagenerators.DefaultTestCoordinate()
	executionLocation := kernel.GeoCoordinate{
		Lat: targetLocation.Latitude() + 0.1,  // Much more significant difference
		Lon: targetLocation.Longitude() + 0.1, // Much more significant difference
	}

	// Get default test data and customize for addresses test
	questData := testdatagenerators.QuestDataWithLocations(targetLocation, executionLocation)

	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questData)
	s.Require().NoError(err)

	// Act - get quest by ID via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.ID().String())
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	foundQuest := httpAssertions.QuestHTTPGetSuccessfully(getResp, err)

	// Verify quest data matches created quest
	singleAssertions := assertions.NewQuestSingleAssertions(s.Assert())
	singleAssertions.QuestHTTPMatchesDomain(foundQuest, createdQuest)
	singleAssertions.QuestHTTPHasValidLocationData(foundQuest)
	singleAssertions.QuestHTTPHasDifferentLocations(foundQuest)
}

func (s *Suite) TestGetQuestByIDHTTPEmptyArrays() {
	ctx := context.Background()

	// Pre-condition - create quest with empty Equipment and Skills arrays
	emptyData := testdatagenerators.EmptyArraysQuestData()

	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, emptyData)
	s.Require().NoError(err)

	// Act - get quest by ID via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.ID().String())
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	foundQuest := httpAssertions.QuestHTTPGetSuccessfully(getResp, err)

	// Verify quest data matches created quest
	singleAssertions := assertions.NewQuestSingleAssertions(s.Assert())
	singleAssertions.QuestHTTPMatchesDomain(foundQuest, createdQuest)

	// Specifically verify that empty arrays are returned as [] and not null
	s.Assert().NotNil(foundQuest.Equipment, "Equipment should not be null")
	s.Assert().NotNil(foundQuest.Skills, "Skills should not be null")
	s.Assert().Len(*foundQuest.Equipment, 0, "Equipment should be empty array")
	s.Assert().Len(*foundQuest.Skills, 0, "Skills should be empty array")
}
