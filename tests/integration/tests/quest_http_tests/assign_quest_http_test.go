package quest_http_tests

import (
	"context"

	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	"quest-manager/tests/integration/mock"
)

// API layer validation tests driven by OpenAPI request validation

func (s *Suite) TestAssignQuestHTTP() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	assignAssertions := assertions.NewQuestAssignAssertions(s.Assert())

	// Pre-condition - create quest directly via handler (faster, no HTTP overhead)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - assign quest via HTTP API (user ID comes from JWT token)
	assignReq := casesteps.AssignQuestHTTPRequest(createdQuest.ID())
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Assert - use helper to eliminate boilerplate
	assignResult := httpAssertions.QuestHTTPAssignedSuccessfully(assignResp, err)

	// Verify assignment result - user ID is taken from mock auth (DefaultUserID)
	expectedUserID := mock.NewAlwaysSuccessAuthClient().DefaultUserID
	assignAssertions.VerifyQuestAssignmentResponse(&assignResult, createdQuest.ID(), expectedUserID)
}

func (s *Suite) TestAssignQuestHTTPInvalidQuestIDFormat() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Test cases with invalid quest ID formats in URL path (OpenAPI path parameter validation)
	testCases := []struct {
		name    string
		questID string
	}{
		{
			name:    "invalid quest ID format",
			questID: "not-a-uuid",
		},
		{
			name:    "partial quest ID",
			questID: "123e4567-e89b-12d3-a456",
		},
		{
			name:    "empty quest ID",
			questID: "",
		},
		{
			name:    "numeric quest ID",
			questID: "12345",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act - send request with invalid quest_id in URL path
			assignReq := casesteps.AssignQuestHTTPRequestWithStringID(tc.questID)
			assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

			// Assert - API layer should reject invalid quest ID format
			httpAssertions.QuestHTTPValidationError(assignResp, err, "questId")
			if tc.questID == "" {
				s.Assert().Contains(assignResp.Body, "missing", "Error should mention missing quest ID")
			}
		})
	}
}

func (s *Suite) TestAssignQuestHTTPPersistence() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	singleAssertions := assertions.NewQuestSingleAssertions(s.Assert())

	// Pre-condition - create quest directly via handler (faster, no HTTP overhead)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - assign quest via HTTP API (user ID comes from JWT token)
	assignReq := casesteps.AssignQuestHTTPRequest(createdQuest.ID())
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Assert assignment using helper
	httpAssertions.QuestHTTPAssignedSuccessfully(assignResp, err)

	// Verify quest is persisted with assignment by retrieving it via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.ID())
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert retrieval and assignment state using helper
	// User ID is taken from mock auth (DefaultUserID)
	expectedUserID := mock.NewAlwaysSuccessAuthClient().DefaultUserID
	retrievedQuest := httpAssertions.QuestHTTPGetSuccessfully(getResp, err)
	singleAssertions.QuestHTTPIsAssignedToUser(retrievedQuest, expectedUserID, createdQuest.ID())
}
