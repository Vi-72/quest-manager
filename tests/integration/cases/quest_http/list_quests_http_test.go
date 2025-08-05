package quest_http

// API LAYER VALIDATION TESTS
// Only tests that correspond to HTTP API layer functionality

import (
	"context"
	"encoding/json"
	"net/http"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestListQuestsHTTP() {
	ctx := context.Background()

	// Pre-condition - create multiple quests via handler (for setup)
	expectedCount := 2
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Act - get list of quests via HTTP API (all quests, without status filter)
	listReq := casesteps.ListQuestsHTTPRequest("")
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	quests := httpAssertions.QuestHTTPListSuccessfully(listResp, err)

	// Verify response
	s.Assert().GreaterOrEqual(len(quests), expectedCount, "Should return at least %d quests", expectedCount)

	// Verify all created quests are in the response
	returnedQuestIDs := make(map[string]bool)
	for _, q := range quests {
		returnedQuestIDs[q.Id] = true
	}

	for _, createdQuest := range createdQuests {
		questID := createdQuest.ID().String()
		s.Assert().True(returnedQuestIDs[questID], "Created quest %s should be in quests list", questID)
	}
}

func (s *Suite) TestListQuestsHTTPEmpty() {
	ctx := context.Background()

	// Act - get list of quests from empty database via HTTP API
	listReq := casesteps.ListQuestsHTTPRequest("")
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, listResp.StatusCode)

	// Parse response
	var quests []servers.Quest
	err = json.Unmarshal([]byte(listResp.Body), &quests)
	s.Require().NoError(err)

	s.Assert().Len(quests, 0, "Should return empty list when no quests exist")
}

func (s *Suite) TestListQuestsHTTPWithValidStatus() {
	ctx := context.Background()

	// Pre-condition - create multiple quests via handler
	expectedCount := 3
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Change quest status for filtering test
	targetStatus := quest.StatusPosted
	_, err = casesteps.ChangeQuestStatusStep(ctx, s.TestDIContainer.ChangeQuestStatusHandler,
		s.TestDIContainer.QuestRepository, createdQuests[0].ID(), targetStatus)
	s.Require().NoError(err)

	// Leave other quests with default StatusCreated status

	// Act - get list of quests filtered by StatusPosted via HTTP API
	listReq := casesteps.ListQuestsHTTPRequest(string(targetStatus))
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, listResp.StatusCode)

	// Parse response
	var quests []servers.Quest
	err = json.Unmarshal([]byte(listResp.Body), &quests)
	s.Require().NoError(err)

	s.Assert().GreaterOrEqual(len(quests), 1, "Should have at least one quest with StatusPosted")

	// Verify all returned quests have the correct status
	for _, q := range quests {
		s.Assert().Equal(string(targetStatus), string(q.Status), "All quests should have StatusPosted")
	}

	// Verify our specific quest is in the list
	foundTargetQuest := false
	targetQuestID := createdQuests[0].ID().String()
	for _, q := range quests {
		if q.Id == targetQuestID {
			foundTargetQuest = true
			break
		}
	}
	s.Assert().True(foundTargetQuest, "Quest with StatusPosted should be in filtered list")
}

func (s *Suite) TestListQuestsHTTPWithEmptyStatus() {
	ctx := context.Background()

	// Pre-condition - create multiple quests via handler
	expectedCount := 2
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Act - get list of quests without status filter via HTTP API
	listReq := casesteps.ListQuestsHTTPRequest("")
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, listResp.StatusCode)

	// Parse response
	var quests []servers.Quest
	err = json.Unmarshal([]byte(listResp.Body), &quests)
	s.Require().NoError(err)

	// Verify response
	s.Assert().GreaterOrEqual(len(quests), expectedCount, "Should return at least %d quests", expectedCount)

	// Verify all created quests are in the response
	returnedQuestIDs := make(map[string]bool)
	for _, q := range quests {
		returnedQuestIDs[q.Id] = true
	}

	for _, createdQuest := range createdQuests {
		questID := createdQuest.ID().String()
		s.Assert().True(returnedQuestIDs[questID], "Created quest %s should be in quests list", questID)
	}
}

func (s *Suite) TestListQuestsHTTPWithInvalidStatus() {
	ctx := context.Background()

	// Pre-condition - create quest to ensure database is not empty
	_, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - try to get list of quests with invalid status via HTTP API
	listReq := casesteps.ListQuestsHTTPRequest("invalid_status_that_does_not_exist")
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert - should return 400 validation error
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, listResp.StatusCode, "Should return 400 for invalid status")

	// Verify error message contains validation details
	s.Assert().Contains(listResp.Body, "validation failed", "Error message should contain validation failure details")
	s.Assert().Contains(listResp.Body, "must be one of", "Error message should mention valid status values")
}
