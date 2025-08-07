package assertions

import (
	"quest-manager/internal/core/domain/model/quest"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"

	"github.com/stretchr/testify/assert"
)

// QuestHandlerAssertions provides reusable assertions for quest handler tests
type QuestHandlerAssertions struct {
	assert *assert.Assertions
}

// NewQuestHandlerAssertions creates a new instance of quest handler assertions
func NewQuestHandlerAssertions(assert *assert.Assertions) *QuestHandlerAssertions {
	return &QuestHandlerAssertions{
		assert: assert,
	}
}

// VerifyQuestFullMatch verifies that a created quest fully matches test data with no error
func (a *QuestHandlerAssertions) VerifyQuestFullMatch(createdQuest quest.Quest, questData testdatagenerators.QuestTestData, err error) {
	// First verify no error occurred
	a.assert.NoError(err, "Quest creation should succeed")

	// Verify basic quest properties
	a.assert.NotEmpty(createdQuest.ID().String(), "Quest should have a valid ID")
	a.assert.Equal(questData.Title, createdQuest.Title, "Quest title should match test data")
	a.assert.Equal(questData.Description, createdQuest.Description, "Quest description should match test data")
	a.assert.Equal(questData.Difficulty, string(createdQuest.Difficulty), "Quest difficulty should match test data")
	a.assert.Equal(questData.Reward, createdQuest.Reward, "Quest reward should match test data")
	a.assert.Equal(questData.DurationMinutes, createdQuest.DurationMinutes, "Quest duration should match test data")
	a.assert.Equal(questData.Creator, createdQuest.Creator, "Quest creator should match test data")

	// Verify arrays (equipment and skills)
	a.assert.Equal(questData.Equipment, createdQuest.Equipment, "Quest equipment should match test data")
	a.assert.Equal(questData.Skills, createdQuest.Skills, "Quest skills should match test data")

	// Verify location coordinates
	a.assert.InDelta(questData.TargetLocation.Latitude(), createdQuest.TargetLocation.Latitude(), 0.001, "Target location latitude should match")
	a.assert.InDelta(questData.TargetLocation.Longitude(), createdQuest.TargetLocation.Longitude(), 0.001, "Target location longitude should match")
	a.assert.InDelta(questData.ExecutionLocation.Latitude(), createdQuest.ExecutionLocation.Latitude(), 0.001, "Execution location latitude should match")
	a.assert.InDelta(questData.ExecutionLocation.Longitude(), createdQuest.ExecutionLocation.Longitude(), 0.001, "Execution location longitude should match")

	// Verify timestamps are set (handler level should set these)
	a.assert.False(createdQuest.CreatedAt.IsZero(), "CreatedAt should be set")
	a.assert.False(createdQuest.UpdatedAt.IsZero(), "UpdatedAt should be set")

	// Verify default status
	a.assert.Equal(quest.StatusCreated, createdQuest.Status, "New quest should have Created status")
}
