package quest_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
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
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	// Parse response
	var foundQuest servers.Quest
	err = json.Unmarshal([]byte(getResp.Body), &foundQuest)
	s.Require().NoError(err)

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
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, getResp.StatusCode, "Should return 404 for non-existent quest ID")

	// Verify error message
	s.Assert().Contains(getResp.Body, "not found", "Error message should indicate quest was not found")
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
	defaultData := testdatagenerators.DefaultQuestData()

	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler,
		"Test Quest with Addresses", "Test Description with Addresses", defaultData.Difficulty,
		defaultData.Reward, defaultData.DurationMinutes, defaultData.Creator,
		targetLocation, executionLocation,
		defaultData.Equipment, defaultData.Skills)
	s.Require().NoError(err)

	// Act - get quest by ID via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.ID().String())
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	// Parse response
	var foundQuest servers.Quest
	err = json.Unmarshal([]byte(getResp.Body), &foundQuest)
	s.Require().NoError(err)

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

	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler,
		emptyData.Title, emptyData.Description, emptyData.Difficulty,
		emptyData.Reward, emptyData.DurationMinutes, emptyData.Creator,
		emptyData.TargetLocation, emptyData.ExecutionLocation,
		emptyData.Equipment, emptyData.Skills)
	s.Require().NoError(err)

	// Act - get quest by ID via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.ID().String())
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	// Parse response
	var foundQuest servers.Quest
	err = json.Unmarshal([]byte(getResp.Body), &foundQuest)
	s.Require().NoError(err)

	// Verify quest data matches created quest
	singleAssertions := assertions.NewQuestSingleAssertions(s.Assert())
	singleAssertions.QuestHTTPMatchesDomain(foundQuest, createdQuest)

	// Specifically verify that empty arrays are returned as [] and not null
	s.Assert().NotNil(foundQuest.Equipment, "Equipment should not be null")
	s.Assert().NotNil(foundQuest.Skills, "Skills should not be null")
	s.Assert().Len(*foundQuest.Equipment, 0, "Equipment should be empty array")
	s.Assert().Len(*foundQuest.Skills, 0, "Skills should be empty array")
}
