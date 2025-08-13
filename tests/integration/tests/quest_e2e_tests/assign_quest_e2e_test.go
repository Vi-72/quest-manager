package quest_e2e_tests

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"

	"github.com/google/uuid"
)

// Test 3: Create quest via handler, assign via API, verify database and events
func (s *E2ESuite) TestCreateThroughHandlerAssignThroughAPI() {
	ctx := context.Background()

	// 1. Create quest through handler (not API) using helper
	questData := testdatagenerators.SimpleQuestData(
		"Handler Created Quest",
		"Quest created via handler for E2E assign test",
		"medium",
		3,
		90,
		testdatagenerators.DefaultTestCoordinate(),
		testdatagenerators.DefaultTestCoordinate(),
	)

	createdQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questData)
	s.Require().NoError(err)

	// Wait for async processing
	time.Sleep(100 * time.Millisecond)

	// 2. Assign quest through API using helper
	userID := uuid.New().String()
	assignReq := casesteps.AssignQuestHTTPRequest(createdQuest.ID().String(), userID)
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Check HTTP response
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusOK, assignResp.StatusCode, "Quest assignment should succeed")

	var assignResult servers.AssignQuestResult
	err = json.Unmarshal([]byte(assignResp.Body), &assignResult)
	s.Require().NoError(err, "Should unmarshal assign response")

	// Wait for async processing
	time.Sleep(100 * time.Millisecond)

	// 3. Verify database data is updated
	updatedQuest, err := s.TestDIContainer.QuestRepository.GetByID(ctx, createdQuest.ID())
	s.Require().NoError(err, "Should retrieve updated quest from database")
	s.Assert().Equal(quest.StatusAssigned, updatedQuest.Status, "Quest status should be assigned")
	s.Assert().NotNil(updatedQuest.Assignee, "Assignee field should be set")
	s.Assert().Equal(userID, *updatedQuest.Assignee, "Assignee ID should match")

	// Verify response data matches DB using assign assertions
	assignAssertions := assertions.NewQuestAssignAssertions(s.Assert())
	assignAssertions.VerifyQuestAssignmentResponse(&assignResult, createdQuest.ID(), userID)

	// 4. Verify assignment events were published
	s.Assert().True(updatedQuest.UpdatedAt.After(updatedQuest.CreatedAt), "Update time should be after creation time")
	s.Assert().Equal(quest.StatusAssigned, updatedQuest.Status, "Status change indicates assign event was processed")
	s.Assert().NotNil(updatedQuest.Assignee, "Assignee field indicates assign event was processed")
	s.Assert().True(updatedQuest.UpdatedAt.After(createdQuest.UpdatedAt), "Updated time change indicates assign event was processed")

	assignedQuests, err := s.TestDIContainer.QuestRepository.FindByAssignee(ctx, userID)
	s.Require().NoError(err, "Should find quests by assignee")
	s.Assert().Len(assignedQuests, 1, "Should find exactly one assigned quest")
	s.Assert().Equal(createdQuest.ID(), assignedQuests[0].ID(), "Found quest should match the assigned quest")

	listAssignedResult, err := casesteps.ListAssignedQuestsStep(ctx, s.TestDIContainer.ListAssignedQuestsHandler, userID)
	s.Require().NoError(err, "Should list assigned quests successfully")
	s.Assert().Len(listAssignedResult, 1, "Should have one assigned quest in list")
	s.Assert().Equal(createdQuest.ID(), listAssignedResult[0].ID(), "Listed quest should match assigned quest")

	// 5. Verify assignment events were saved in the events table
	e2eAssertions := assertions.NewQuestE2EAssertions(s.Assert(), s.TestDIContainer.EventStorage)
	e2eAssertions.VerifyQuestAssignmentEvents(ctx, createdQuest.ID(), createdQuest.CreatedAt)
}
