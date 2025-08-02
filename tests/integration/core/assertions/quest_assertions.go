package assertions

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/tests/integration/core/storage"
)

// QuestAssertions provides methods for asserting quests in tests
type QuestAssertions struct {
	assert       *assert.Assertions
	require      *require.Assertions
	storage      *storage.QuestStorage
	eventStorage *storage.EventStorage
}

// QuestExists verifies that quest exists in database
func (a *QuestAssertions) QuestExists(ctx context.Context, questID uuid.UUID) {
	q, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Quest should exist in database")
	a.require.NotNil(q, "Quest should not be nil")
	a.assert.Equal(questID.String(), q.ID, "Quest ID should match")
}

// QuestNotExists verifies that quest does not exist in database
func (a *QuestAssertions) QuestNotExists(ctx context.Context, questID uuid.UUID) {
	_, err := a.storage.GetQuestByID(ctx, questID)
	a.assert.Error(err, "Quest should not exist in database")
}

// QuestHasStatus verifies quest status
func (a *QuestAssertions) QuestHasStatus(ctx context.Context, questID uuid.UUID, expectedStatus quest.Status) {
	q, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Failed to get quest from database")
	a.assert.Equal(string(expectedStatus), q.Status, "Quest status should match expected")
}

// QuestHasAssignee verifies quest assignee
func (a *QuestAssertions) QuestHasAssignee(ctx context.Context, questID uuid.UUID, expectedAssignee *string) {
	q, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Failed to get quest from database")

	if expectedAssignee == nil {
		a.assert.Nil(q.Assignee, "Quest should not have assignee")
	} else {
		a.require.NotNil(q.Assignee, "Quest should have assignee")
		a.assert.Equal(*expectedAssignee, *q.Assignee, "Quest assignee should match expected")
	}
}

// QuestHasCreator verifies quest creator
func (a *QuestAssertions) QuestHasCreator(ctx context.Context, questID uuid.UUID, expectedCreator string) {
	q, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Failed to get quest from database")
	a.assert.Equal(expectedCreator, q.Creator, "Quest creator should match expected")
}

// QuestCountEquals verifies total quest count
func (a *QuestAssertions) QuestCountEquals(ctx context.Context, expectedCount int64) {
	count, err := a.storage.CountQuests(ctx)
	a.require.NoError(err, "Failed to count quests")
	a.assert.Equal(expectedCount, count, "Quest count should match expected")
}

// QuestCountByStatusEquals verifies quest count by status
func (a *QuestAssertions) QuestCountByStatusEquals(ctx context.Context, status quest.Status, expectedCount int64) {
	count, err := a.storage.CountQuestsByStatus(ctx, status)
	a.require.NoError(err, "Failed to count quests by status")
	a.assert.Equal(expectedCount, count, "Quest count by status should match expected")
}

// QuestHasValidTimestamps verifies that quest timestamps are correct
func (a *QuestAssertions) QuestHasValidTimestamps(ctx context.Context, questID uuid.UUID) {
	q, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Failed to get quest from database")

	a.assert.False(q.CreatedAt.IsZero(), "CreatedAt should not be zero")
	a.assert.False(q.UpdatedAt.IsZero(), "UpdatedAt should not be zero")
	a.assert.True(q.UpdatedAt.After(q.CreatedAt) || q.UpdatedAt.Equal(q.CreatedAt),
		"UpdatedAt should be after or equal to CreatedAt")
}

// EventExists verifies that event exists
func (a *QuestAssertions) EventExists(ctx context.Context, eventType string, aggregateID uuid.UUID) {
	events, err := a.eventStorage.GetEventsByAggregateID(ctx, aggregateID)
	a.require.NoError(err, "Failed to get events")

	found := false
	for _, event := range events {
		if event.EventType == eventType {
			found = true
			break
		}
	}

	a.assert.True(found, "Event of type %s should exist for aggregate %s", eventType, aggregateID)
}

// EventCountEquals verifies event count for aggregate
func (a *QuestAssertions) EventCountEquals(ctx context.Context, aggregateID uuid.UUID, expectedCount int64) {
	count, err := a.eventStorage.CountEventsByAggregateID(ctx, aggregateID)
	a.require.NoError(err, "Failed to count events")
	a.assert.Equal(expectedCount, count, "Event count should match expected")
}

// EventCountByTypeEquals verifies event count by type
func (a *QuestAssertions) EventCountByTypeEquals(ctx context.Context, eventType string, expectedCount int64) {
	count, err := a.eventStorage.CountEventsByType(ctx, eventType)
	a.require.NoError(err, "Failed to count events by type")
	a.assert.Equal(expectedCount, count, "Event count by type should match expected")
}

// WaitForEvents waits for events to appear (for asynchronous operations)
func (a *QuestAssertions) WaitForEvents(ctx context.Context, aggregateID uuid.UUID, expectedCount int64, timeout time.Duration) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		count, err := a.eventStorage.CountEventsByAggregateID(ctx, aggregateID)
		a.require.NoError(err, "Failed to count events")

		if count >= expectedCount {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	// Final check
	count, err := a.eventStorage.CountEventsByAggregateID(ctx, aggregateID)
	a.require.NoError(err, "Failed to count events")
	a.assert.GreaterOrEqual(count, expectedCount, "Should have at least %d events after timeout", expectedCount)
}
