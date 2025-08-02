package assertions

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/quest"
)

type QuestSingleAssertions struct {
	assert *assert.Assertions
}

func NewQuestSingleAssertions(a *assert.Assertions) *QuestSingleAssertions {
	return &QuestSingleAssertions{assert: a}
}

// QuestMatchesCreated verifies that found quest completely matches the created quest
func (a *QuestSingleAssertions) QuestMatchesCreated(found, created quest.Quest) {
	a.assert.Equal(created.ID(), found.ID(), "Quest ID should match")
	a.assert.Equal(created.Title, found.Title, "Quest title should match")
	a.assert.Equal(created.Description, found.Description, "Quest description should match")
	a.assert.Equal(string(created.Difficulty), string(found.Difficulty), "Quest difficulty should match")
	a.assert.Equal(created.Reward, found.Reward, "Quest reward should match")
	a.assert.Equal(created.DurationMinutes, found.DurationMinutes, "Quest duration should match")
	a.assert.Equal(created.Creator, found.Creator, "Quest creator should match")
	a.assert.Equal(created.Status, found.Status, "Quest status should match")

	a.assert.False(found.CreatedAt.IsZero(), "CreatedAt should not be zero")
	a.assert.False(found.UpdatedAt.IsZero(), "UpdatedAt should not be zero")
	a.assert.True(found.UpdatedAt.After(found.CreatedAt) || found.UpdatedAt.Equal(found.CreatedAt),
		"UpdatedAt should be after or equal to CreatedAt")
}

// QuestHasValidLocationData verifies that quest has valid location IDs and addresses
func (a *QuestSingleAssertions) QuestHasValidLocationData(q quest.Quest) {
	a.assert.NotNil(q.TargetLocationID, "Quest should have target location ID")
	a.assert.NotEqual(uuid.Nil, *q.TargetLocationID, "Target location ID should not be nil UUID")

	a.assert.NotNil(q.ExecutionLocationID, "Quest should have execution location ID")
	a.assert.NotEqual(uuid.Nil, *q.ExecutionLocationID, "Execution location ID should not be nil UUID")

	a.assert.NotNil(q.TargetAddress, "Quest should have target address")
	a.assert.NotEmpty(*q.TargetAddress, "Target address should not be empty")

	a.assert.NotNil(q.ExecutionAddress, "Quest should have execution address")
	a.assert.NotEmpty(*q.ExecutionAddress, "Execution address should not be empty")
}
