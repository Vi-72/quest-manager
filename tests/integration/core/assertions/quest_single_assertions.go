package assertions

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
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

// QuestHTTPMatchesDomain verifies that HTTP quest response matches domain quest
func (a *QuestSingleAssertions) QuestHTTPMatchesDomain(httpQuest servers.Quest, domainQuest quest.Quest) {
	a.assert.Equal(domainQuest.ID().String(), httpQuest.Id, "Quest ID should match")
	a.assert.Equal(domainQuest.Title, httpQuest.Title, "Quest title should match")
	a.assert.Equal(domainQuest.Description, httpQuest.Description, "Quest description should match")
	a.assert.Equal(string(domainQuest.Difficulty), string(httpQuest.Difficulty), "Quest difficulty should match")
	a.assert.Equal(domainQuest.Reward, httpQuest.Reward, "Quest reward should match")
	a.assert.Equal(domainQuest.DurationMinutes, httpQuest.DurationMinutes, "Quest duration should match")
	a.assert.Equal(domainQuest.Creator, httpQuest.Creator, "Quest creator should match")
	a.assert.Equal(string(domainQuest.Status), string(httpQuest.Status), "Quest status should match")
}

// QuestHTTPHasValidLocationData verifies that HTTP quest has valid location IDs and coordinates
func (a *QuestSingleAssertions) QuestHTTPHasValidLocationData(q servers.Quest) {
	a.assert.NotNil(q.TargetLocationId, "Quest should have target location ID")
	a.assert.NotEmpty(*q.TargetLocationId, "Target location ID should not be empty")
	a.assert.NotNil(q.ExecutionLocationId, "Quest should have execution location ID")
	a.assert.NotEmpty(*q.ExecutionLocationId, "Execution location ID should not be empty")

	a.assert.NotZero(q.TargetLocation.Latitude, "Target location should have latitude")
	a.assert.NotZero(q.TargetLocation.Longitude, "Target location should have longitude")
	a.assert.NotZero(q.ExecutionLocation.Latitude, "Execution location should have latitude")
	a.assert.NotZero(q.ExecutionLocation.Longitude, "Execution location should have longitude")
}

// QuestHTTPHasDifferentLocations verifies that target and execution locations are different
func (a *QuestSingleAssertions) QuestHTTPHasDifferentLocations(q servers.Quest) {
	a.assert.NotNil(q.TargetLocationId, "Quest should have target location ID")
	a.assert.NotNil(q.ExecutionLocationId, "Quest should have execution location ID")
	a.assert.NotEqual(*q.TargetLocationId, *q.ExecutionLocationId, "Target and execution location IDs should be different")

	a.assert.NotEqual(q.TargetLocation.Latitude, q.ExecutionLocation.Latitude, "Target and execution locations should have different latitudes")
	a.assert.NotEqual(q.TargetLocation.Longitude, q.ExecutionLocation.Longitude, "Target and execution locations should have different longitudes")
}
