package quest_http

import (
	"context"
	"encoding/json"
	"net/http"

	"quest-manager/internal/generated/servers"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateChangeQuestStatusRequest function

func (s *Suite) TestChangeQuestStatusHTTPValidation() {
	ctx := context.Background()

	// Pre-condition - create quest via handler (for setup)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - valid status change
	statusRequest := &servers.ChangeStatusRequest{
		Status: servers.QuestStatusPosted,
	}

	changeReq := casesteps.ChangeQuestStatusHTTPRequest(createdQuest.ID().String(), statusRequest)
	changeResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, changeReq)

	// Assert - successful change (200 OK for status change, not 201)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, changeResp.StatusCode)

	// Parse response manually since we don't have a generic 200 OK assertion yet
	var result servers.ChangeQuestStatusResult
	parseErr := json.Unmarshal([]byte(changeResp.Body), &result)
	s.Require().NoError(parseErr)

	s.Assert().Equal(createdQuest.ID().String(), result.Id)
	s.Assert().Equal(string(servers.QuestStatusPosted), string(result.Status))
}

func (s *Suite) TestChangeQuestStatusHTTPMissingBody() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - send request with empty object to test ValidateBody function
	emptyBodyRequest := map[string]interface{}{} // Empty object (missing required fields)
	changeReq := casesteps.ChangeQuestStatusHTTPRequest(createdQuest.ID().String(), emptyBodyRequest)
	changeResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, changeReq)

	// Assert - API layer should reject incomplete body
	httpAssertions.QuestHTTPValidationError(changeResp, err, "status")
}

func (s *Suite) TestChangeQuestStatusHTTPInvalidUUID() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Act - send request with invalid UUID format
	statusRequest := &servers.ChangeStatusRequest{
		Status: servers.QuestStatusPosted,
	}

	changeReq := casesteps.ChangeQuestStatusHTTPRequest("invalid-uuid-format", statusRequest)
	changeResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, changeReq)

	// Assert - API layer should reject invalid UUID
	httpAssertions.QuestHTTPValidationError(changeResp, err, "questId")
}

func (s *Suite) TestChangeQuestStatusHTTPEmptyStatus() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - send request with empty status (ValidateNotEmpty function test)
	emptyStatusRequest := map[string]interface{}{
		"status": "",
	}
	changeReq := casesteps.ChangeQuestStatusHTTPRequest(createdQuest.ID().String(), emptyStatusRequest)
	changeResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, changeReq)

	// Assert - API layer should reject empty status
	httpAssertions.QuestHTTPValidationError(changeResp, err, "status")
}

func (s *Suite) TestChangeQuestStatusHTTPMalformedJSON() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - create quest
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - send malformed JSON
	changeReq := casesteps.CreateMalformedJSONRequest("PATCH", "/api/v1/quests/"+createdQuest.ID().String()+"/status")
	changeResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, changeReq)

	// Assert - should return 400 for malformed JSON
	httpAssertions.QuestHTTPErrorResponse(changeResp, err, http.StatusBadRequest, "")
}
