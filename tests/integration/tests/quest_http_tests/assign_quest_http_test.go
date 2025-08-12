package quest_http_tests

import (
	"context"
	"encoding/json"
	"net/http"

	"quest-manager/internal/generated/servers"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"

	"github.com/google/uuid"
)

// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateAssignQuestRequest function

func (s *Suite) TestAssignQuestHTTP() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	assignAssertions := assertions.NewQuestAssignAssertions(s.Assert())

	// Pre-condition - create quest
	questRequest := testdatagenerators.RandomCreateQuestRequest()
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Act - assign quest via HTTP API
	userID := "123e4567-e89b-12d3-a456-426614174000" // Valid UUID
	assignReq := casesteps.AssignQuestHTTPRequest(createdQuest.Id, userID)
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, assignResp.StatusCode)

	// Parse response
	var assignResult servers.AssignQuestResult
	parseErr := json.Unmarshal([]byte(assignResp.Body), &assignResult)
	s.Require().NoError(parseErr, "Response should be valid JSON")

	// Verify assignment result
	questID, parseErr := uuid.Parse(createdQuest.Id)
	s.Require().NoError(parseErr, "Created quest ID should be valid UUID")
	assignAssertions.VerifyQuestAssignmentResponse(&assignResult, questID, userID)
}

// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateAssignQuestRequest function

func (s *Suite) TestAssignQuestHTTPMissingRequiredFields() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest
	questRequest := testdatagenerators.RandomCreateQuestRequest()
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Act - send request with empty JSON body to test ValidateBody function
	emptyBodyRequest := map[string]interface{}{} // Empty object

	assignReq := casesteps.HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests/" + createdQuest.Id + "/assign",
		Body:        emptyBodyRequest,
		ContentType: "application/json",
	}
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Assert - API layer should reject incomplete body
	httpAssertions.QuestHTTPValidationError(assignResp, err, "")
}

func (s *Suite) TestAssignQuestHTTPMissingUserID() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest
	questRequest := testdatagenerators.RandomCreateQuestRequest()
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Act - send request without user_id field
	requestBody := map[string]interface{}{
		// user_id is missing
	}

	assignReq := casesteps.HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests/" + createdQuest.Id + "/assign",
		Body:        requestBody,
		ContentType: "application/json",
	}
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Assert - API layer should reject missing user_id
	httpAssertions.QuestHTTPValidationError(assignResp, err, "userId")
}

func (s *Suite) TestAssignQuestHTTPEmptyUserID() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest
	questRequest := testdatagenerators.RandomCreateQuestRequest()
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Act - send request with empty user_id
	requestBody := map[string]interface{}{
		"user_id": "", // Empty user_id - ValidateUUID should catch this
	}

	assignReq := casesteps.HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests/" + createdQuest.Id + "/assign",
		Body:        requestBody,
		ContentType: "application/json",
	}
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Assert - API layer should reject empty user_id
	httpAssertions.QuestHTTPValidationError(assignResp, err, "userId")
}

func (s *Suite) TestAssignQuestHTTPInvalidUserIDFormat() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest
	questRequest := testdatagenerators.RandomCreateQuestRequest()
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Test cases with invalid UUID formats (ValidateUUID function tests)
	testCases := []struct {
		name   string
		userID string
	}{
		{
			name:   "invalid UUID format",
			userID: "not-a-uuid",
		},
		{
			name:   "partial UUID",
			userID: "123e4567-e89b-12d3-a456",
		},
		{
			name:   "too long string",
			userID: "123e4567-e89b-12d3-a456-426614174000-extra",
		},
		{
			name:   "numeric string",
			userID: "12345",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act - send request with invalid user_id format
			requestBody := map[string]interface{}{
				"user_id": tc.userID,
			}

			assignReq := casesteps.HTTPRequest{
				Method:      "POST",
				URL:         "/api/v1/quests/" + createdQuest.Id + "/assign",
				Body:        requestBody,
				ContentType: "application/json",
			}
			assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

			// Assert - API layer should reject invalid UUID format
			httpAssertions.QuestHTTPValidationError(assignResp, err, "userId")
		})
	}
}

func (s *Suite) TestAssignQuestHTTPInvalidQuestIDFormat() {
	ctx := context.Background()

	// Test cases with invalid quest ID formats in URL path (ValidateUUID function tests)
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
			requestBody := map[string]interface{}{
				"user_id": "123e4567-e89b-12d3-a456-426614174000", // Valid UUID
			}

			assignReq := casesteps.HTTPRequest{
				Method:      "POST",
				URL:         "/api/v1/quests/" + tc.questID + "/assign",
				Body:        requestBody,
				ContentType: "application/json",
			}
			assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

			// Assert - API layer should reject invalid quest ID format
			s.Require().NoError(err, "HTTP request should not fail")
			s.Require().Equal(http.StatusBadRequest, assignResp.StatusCode, "Should return 400 for invalid quest ID")
			// Framework level validation message may be different from application level
			if tc.questID == "" {
				s.Assert().Contains(assignResp.Body, "empty", "Error should mention empty quest ID")
			} else {
				s.Assert().Contains(assignResp.Body, "questId", "Error should mention quest ID field")
			}
		})
	}
}

func (s *Suite) TestAssignQuestHTTPMalformedJSON() {
	ctx := context.Background()

	// Pre-condition - create quest
	questRequest := testdatagenerators.RandomCreateQuestRequest()
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Act - send malformed JSON
	malformedRequest := casesteps.HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests/" + createdQuest.Id + "/assign",
		Body:        `{"user_id": "invalid-json", }`, // Malformed JSON
		ContentType: "application/json",
	}

	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, malformedRequest)

	// Assert - API layer should reject malformed JSON
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, assignResp.StatusCode, "Should return 400 for malformed JSON")
}

func (s *Suite) TestAssignQuestHTTPPersistence() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest
	questRequest := testdatagenerators.RandomCreateQuestRequest()
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Act - assign quest via HTTP API
	userID := "123e4567-e89b-12d3-a456-426614174000"
	assignReq := casesteps.AssignQuestHTTPRequest(createdQuest.Id, userID)
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Assert assignment
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, assignResp.StatusCode)

	// Verify quest is persisted with assignment by retrieving it via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.Id)
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert retrieval
	retrievedQuest := httpAssertions.QuestHTTPGetSuccessfully(getResp, err)

	// Verify quest is assigned
	s.Assert().Equal(createdQuest.Id, retrievedQuest.Id)
	s.Assert().NotNil(retrievedQuest.Assignee, "Quest should have assignee")
	s.Assert().Equal(userID, *retrievedQuest.Assignee)
	s.Assert().Equal(servers.QuestStatusAssigned, retrievedQuest.Status)
}
