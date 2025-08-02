package assertions

import (
	"quest-manager/internal/core/domain/model/quest"

	"github.com/stretchr/testify/assert"
)

// QuestListAssertions provides methods for asserting quest lists in tests
type QuestListAssertions struct {
	assert *assert.Assertions
}

// NewQuestListAssertions creates a new QuestListAssertions
func NewQuestListAssertions(a *assert.Assertions) *QuestListAssertions {
	return &QuestListAssertions{
		assert: a,
	}
}

// QuestsContainAllCreated verifies that all created quests are present in the retrieved list
func (a *QuestListAssertions) QuestsContainAllCreated(createdQuests []quest.Quest, retrievedQuests []quest.Quest) {
	// Create map from retrieved quest IDs for fast lookup
	retrievedQuestIDs := make(map[string]bool)
	for _, q := range retrievedQuests {
		retrievedQuestIDs[q.ID().String()] = true
	}

	// Verify that each created quest is present in the retrieved list
	for _, createdQuest := range createdQuests {
		a.assert.Contains(retrievedQuestIDs, createdQuest.ID().String(),
			"Created quest with ID %s should be in the retrieved list", createdQuest.ID().String())
	}
}

// QuestsHaveMinimumCount verifies that the list contains at least the minimum number of quests
func (a *QuestListAssertions) QuestsHaveMinimumCount(quests []quest.Quest, expectedMinimum int) {
	a.assert.GreaterOrEqual(len(quests), expectedMinimum,
		"Should have at least %d quests, but got %d", expectedMinimum, len(quests))
}

// QuestsAllHaveStatus verifies that all quests in the list have the specified status
func (a *QuestListAssertions) QuestsAllHaveStatus(quests []quest.Quest, expectedStatus quest.Status) {
	for i, q := range quests {
		a.assert.Equal(expectedStatus, q.Status,
			"Quest at index %d should have status %s, but got %s", i, expectedStatus, q.Status)
	}
}

// QuestWithIDExists verifies that a quest with the specified ID is present in the list
func (a *QuestListAssertions) QuestWithIDExists(quests []quest.Quest, questID string) {
	found := false
	for _, q := range quests {
		if q.ID().String() == questID {
			found = true
			break
		}
	}
	a.assert.True(found, "Quest with ID %s should be present in the list", questID)
}

// QuestWithIDNotExists verifies that a quest with the specified ID is NOT present in the list
func (a *QuestListAssertions) QuestWithIDNotExists(quests []quest.Quest, questID string) {
	found := false
	for _, q := range quests {
		if q.ID().String() == questID {
			found = true
			break
		}
	}
	a.assert.False(found, "Quest with ID %s should NOT be present in the list", questID)
}
