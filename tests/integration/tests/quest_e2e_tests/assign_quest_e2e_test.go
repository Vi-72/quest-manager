package quest_e2e_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"

	"github.com/google/uuid"
)

// Test 3: Create quest via handler, assign via API, verify database and events
func (s *E2ESuite) TestCreateThroughHandlerAssignThroughAPI() {
	ctx := context.Background()

	// 1. Create quest through handler (not API)
	questData := testdatagenerators.SimpleQuestData(
		"Handler Created Quest",
		"Quest created via handler for E2E assign test",
		"medium",
		3,
		90,
		testdatagenerators.DefaultTestCoordinate(),
		testdatagenerators.DefaultTestCoordinate(),
	)

	createCmd := commands.CreateQuestCommand{
		Title:             questData.Title,
		Description:       questData.Description,
		Difficulty:        questData.Difficulty,
		Reward:            questData.Reward,
		DurationMinutes:   questData.DurationMinutes,
		Creator:           questData.Creator,
		TargetLocation:    questData.TargetLocation,
		ExecutionLocation: questData.ExecutionLocation,
		Equipment:         questData.Equipment,
		Skills:            questData.Skills,
	}

	createdQuest, err := s.TestDIContainer.CreateQuestHandler.Handle(ctx, createCmd)
	s.Require().NoError(err, "Quest creation through handler should succeed")
	s.Require().NotNil(createdQuest, "Created quest should not be nil")

	// Wait for async processing
	time.Sleep(100 * time.Millisecond)

	// 2. Assign quest through API
	userID := uuid.New().String()
	assignRequest := servers.AssignQuestRequest{
		UserId: userID,
	}

	requestBody, err := json.Marshal(assignRequest)
	s.Require().NoError(err)

	// Make assign request through HTTP API
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/quests/%s/assign", createdQuest.ID()), bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	s.TestDIContainer.HTTPRouter.ServeHTTP(recorder, req)

	// Check HTTP response
	s.Assert().Equal(http.StatusOK, recorder.Code, "Quest assignment should succeed")

	var assignResult servers.AssignQuestResult
	err = json.Unmarshal(recorder.Body.Bytes(), &assignResult)
	s.Require().NoError(err, "Should unmarshal assign response")

	// Wait for async processing
	time.Sleep(100 * time.Millisecond)

	// 3. Verify database data is updated
	updatedQuest, err := s.TestDIContainer.QuestRepository.GetByID(ctx, createdQuest.ID())
	s.Require().NoError(err, "Should retrieve updated quest from database")
	s.Assert().Equal(quest.StatusAssigned, updatedQuest.Status, "Quest status should be assigned")
	s.Assert().NotNil(updatedQuest.Assignee, "Assignee field should be set")
	s.Assert().Equal(userID, *updatedQuest.Assignee, "Assignee ID should match")

	// Verify response data matches DB
	s.Assert().Equal(createdQuest.ID().String(), assignResult.Id, "Response ID should match quest ID")
	s.Assert().Equal(userID, assignResult.Assignee, "Response assignee should match user ID")
	s.Assert().Equal(servers.QuestStatusAssigned, assignResult.Status, "Response status should be assigned")

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
	questEvents, err := s.TestDIContainer.EventStorage.GetEventsByAggregateID(ctx, createdQuest.ID())
	s.Require().NoError(err, "Should retrieve events for the assigned quest")
	s.Assert().GreaterOrEqual(len(questEvents), 2, "Should have at least 2 events (created + assigned)")

	// Check that quest.assigned event exists
	var questAssignedEventFound bool
	var questStatusChangedEventFound bool
	for _, event := range questEvents {
		if event.EventType == "quest.assigned" {
			questAssignedEventFound = true
			s.Assert().Equal(createdQuest.ID().String(), event.AggregateID, "Assign event aggregate ID should match quest ID")
			s.Assert().NotEmpty(event.Data, "Assign event data should not be empty")
			s.Assert().True(event.CreatedAt.After(createdQuest.CreatedAt), "Assign event should be after quest creation")
		}
		if event.EventType == "quest.status_changed" {
			questStatusChangedEventFound = true
			s.Assert().Equal(createdQuest.ID().String(), event.AggregateID, "Status change event aggregate ID should match quest ID")
			s.Assert().NotEmpty(event.Data, "Status change event data should not be empty")
			s.Assert().True(event.CreatedAt.After(createdQuest.CreatedAt), "Status change event should be after quest creation")
		}
	}
	s.Assert().True(questAssignedEventFound, "Should find quest.assigned event in events table")
	s.Assert().True(questStatusChangedEventFound, "Should find quest.status_changed event in events table")
}
