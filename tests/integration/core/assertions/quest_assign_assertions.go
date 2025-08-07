package assertions

import (
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// QuestAssignAssertions provides reusable assertions for quest assignment operations
type QuestAssignAssertions struct {
	assert *assert.Assertions
}

// NewQuestAssignAssertions creates a new instance of quest assign assertions
func NewQuestAssignAssertions(assert *assert.Assertions) *QuestAssignAssertions {
	return &QuestAssignAssertions{
		assert: assert,
	}
}

// VerifyQuestAssignedSuccessfully verifies that quest assignment was successful
func (a *QuestAssignAssertions) VerifyQuestAssignedSuccessfully(
	err error,
	originalQuest quest.Quest,
	assignResult commands.AssignQuestResult,
	userID string,
) {
	a.assert.NoError(err, "Quest assignment should succeed")
	a.assert.Equal(originalQuest.ID(), assignResult.ID, "Assigned quest ID should match original quest")
	a.assert.Equal(userID, assignResult.Assignee, "Quest assignee should match provided user ID")
	a.assert.Equal(string(quest.StatusAssigned), assignResult.Status, "Quest status should be assigned")
}

// VerifyQuestAssignmentResponse verifies HTTP assignment response
func (a *QuestAssignAssertions) VerifyQuestAssignmentResponse(
	response *servers.AssignQuestResult,
	originalQuestID uuid.UUID,
	userID string,
) {
	a.assert.Equal(originalQuestID.String(), response.Id, "HTTP response quest ID should match original")
	a.assert.Equal(userID, response.Assignee, "HTTP response assignee should match user ID")
	a.assert.Equal(servers.QuestStatusAssigned, response.Status, "HTTP response status should be assigned")
}
