package quest_http_tests

// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateCreateQuestRequest function

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestListAssignedQuestsHTTP() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	listAssertions := assertions.NewQuestListAssertions(s.Assert())

	// Pre-condition - create quests via handler and assign them to a specific user
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

	// Assert using helpers to eliminate boilerplate
	assignedQuests := httpAssertions.QuestHTTPListSuccessfully(listResp, err)
	listAssertions.QuestListHTTPAllAssignedToUser(assignedQuests, testUserID, expectedCount)
	listAssertions.QuestListHTTPContainsAllCreated(assignedQuests, createdQuests)
}

func (s *Suite) TestListAssignedQuestsHTTPEmpty() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - use a user ID that has no assigned quests
	nonExistentUserID := uuid.New().String() // Generate new UUID with no quests

	// Act - get list of assigned quests for user with no assignments via HTTP API
	listReq := casesteps.ListAssignedQuestsHTTPRequest(nonExistentUserID)
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert using helper
	assignedQuests := httpAssertions.QuestHTTPListSuccessfully(listResp, err)
	s.Assert().Len(assignedQuests, 0, "Should return empty list for user with no assigned quests")
}

func (s *Suite) TestListAssignedQuestsHTTPInvalidUserID() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Act - try to get assigned quests with invalid user ID via HTTP API
	listReq := casesteps.ListAssignedQuestsHTTPRequest("InvalidUserID")
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert using helper
	httpAssertions.QuestHTTPErrorResponse(listResp, err, http.StatusBadRequest, "validation failed")
	s.Assert().Contains(listResp.Body, "user_id", "Error message should mention user_id field")
	s.Assert().Contains(listResp.Body, "UUID", "Error message should mention UUID format requirement")
}
