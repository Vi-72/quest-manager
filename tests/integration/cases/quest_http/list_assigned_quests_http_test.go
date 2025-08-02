package quest_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"quest-manager/internal/generated/servers"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestListAssignedQuestsHTTP() {
	ctx := context.Background()

	// Arrange - create quests via handler and assign them to a specific user
	testUserID := uuid.New().String() // Generate new UUID
	expectedCount := 2

	// Create quests via handler
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Assign all created quests to the test user
	for _, quest := range createdQuests {
		_, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, quest.ID(), testUserID)
		s.Require().NoError(err)
	}

	// Act - get list of assigned quests via HTTP API
	listReq := casesteps.ListAssignedQuestsHTTPRequest(testUserID)
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, listResp.StatusCode)

	// Parse response
	var assignedQuests []servers.Quest
	err = json.Unmarshal([]byte(listResp.Body), &assignedQuests)
	s.Require().NoError(err)

	// Verify response
	s.Assert().Len(assignedQuests, expectedCount, "Should return %d assigned quests", expectedCount)

	// Verify all returned quests are assigned to the correct user
	returnedQuestIDs := make(map[string]bool)
	for _, quest := range assignedQuests {
		s.Assert().NotNil(quest.Assignee, "Quest should have an assignee")
		s.Assert().Equal(testUserID, *quest.Assignee, "Quest should be assigned to the test user")
		s.Assert().Equal(servers.QuestStatusAssigned, quest.Status, "Quest should have 'assigned' status")
		returnedQuestIDs[quest.Id] = true
	}

	// Verify all created quests are in the response
	for _, createdQuest := range createdQuests {
		questID := createdQuest.ID().String()
		s.Assert().True(returnedQuestIDs[questID], "Created quest %s should be in assigned quests list", questID)
	}
}

func (s *Suite) TestListAssignedQuestsHTTPEmpty() {
	ctx := context.Background()

	// Arrange - use a user ID that has no assigned quests
	nonExistentUserID := uuid.New().String() // Generate new UUID with no quests

	// Act - get list of assigned quests for user with no assignments via HTTP API
	listReq := casesteps.ListAssignedQuestsHTTPRequest(nonExistentUserID)
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, listResp.StatusCode)

	// Parse response
	var assignedQuests []servers.Quest
	err = json.Unmarshal([]byte(listResp.Body), &assignedQuests)
	s.Require().NoError(err)

	s.Assert().Len(assignedQuests, 0, "Should return empty list for user with no assigned quests")
}

func (s *Suite) TestListAssignedQuestsHTTPInvalidUserID() {
	ctx := context.Background()

	// Act - try to get assigned quests with invalid user ID via HTTP API
	listReq := casesteps.ListAssignedQuestsHTTPRequest("InvalidUserID")
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert - should return 400 error for invalid user ID
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, listResp.StatusCode, "Should return 400 for invalid user ID")

	// Verify error message contains validation details
	s.Assert().Contains(listResp.Body, "validation failed", "Error message should contain validation failure details")
	s.Assert().Contains(listResp.Body, "user_id", "Error message should mention user_id field")
	s.Assert().Contains(listResp.Body, "UUID", "Error message should mention UUID format requirement")
}
