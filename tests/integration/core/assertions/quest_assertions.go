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

// QuestAssertions предоставляет методы для проверки квестов в тестах
type QuestAssertions struct {
	assert       *assert.Assertions
	require      *require.Assertions
	storage      *storage.QuestStorage
	eventStorage *storage.EventStorage
}

// NewQuestAssertions создает новый QuestAssertions
func NewQuestAssertions(a *assert.Assertions, r *require.Assertions, s *storage.QuestStorage, es *storage.EventStorage) *QuestAssertions {
	return &QuestAssertions{
		assert:       a,
		require:      r,
		storage:      s,
		eventStorage: es,
	}
}

// QuestExists проверяет что квест существует в базе данных
func (a *QuestAssertions) QuestExists(ctx context.Context, questID uuid.UUID) {
	quest, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Quest should exist in database")
	a.require.NotNil(quest, "Quest should not be nil")
	a.assert.Equal(questID.String(), quest.ID, "Quest ID should match")
}

// QuestNotExists проверяет что квест не существует в базе данных
func (a *QuestAssertions) QuestNotExists(ctx context.Context, questID uuid.UUID) {
	_, err := a.storage.GetQuestByID(ctx, questID)
	a.assert.Error(err, "Quest should not exist in database")
}

// QuestHasStatus проверяет статус квеста
func (a *QuestAssertions) QuestHasStatus(ctx context.Context, questID uuid.UUID, expectedStatus quest.Status) {
	quest, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Failed to get quest from database")
	a.assert.Equal(string(expectedStatus), quest.Status, "Quest status should match expected")
}

// QuestHasAssignee проверяет назначенного пользователя квеста
func (a *QuestAssertions) QuestHasAssignee(ctx context.Context, questID uuid.UUID, expectedAssignee *string) {
	quest, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Failed to get quest from database")

	if expectedAssignee == nil {
		a.assert.Nil(quest.Assignee, "Quest should not have assignee")
	} else {
		a.require.NotNil(quest.Assignee, "Quest should have assignee")
		a.assert.Equal(*expectedAssignee, *quest.Assignee, "Quest assignee should match expected")
	}
}

// QuestHasCreator проверяет создателя квеста
func (a *QuestAssertions) QuestHasCreator(ctx context.Context, questID uuid.UUID, expectedCreator string) {
	quest, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Failed to get quest from database")
	a.assert.Equal(expectedCreator, quest.Creator, "Quest creator should match expected")
}

// QuestCountEquals проверяет общее количество квестов
func (a *QuestAssertions) QuestCountEquals(ctx context.Context, expectedCount int64) {
	count, err := a.storage.CountQuests(ctx)
	a.require.NoError(err, "Failed to count quests")
	a.assert.Equal(expectedCount, count, "Quest count should match expected")
}

// QuestCountByStatusEquals проверяет количество квестов по статусу
func (a *QuestAssertions) QuestCountByStatusEquals(ctx context.Context, status quest.Status, expectedCount int64) {
	count, err := a.storage.CountQuestsByStatus(ctx, status)
	a.require.NoError(err, "Failed to count quests by status")
	a.assert.Equal(expectedCount, count, "Quest count by status should match expected")
}

// QuestHasValidTimestamps проверяет что временные метки квеста корректны
func (a *QuestAssertions) QuestHasValidTimestamps(ctx context.Context, questID uuid.UUID) {
	quest, err := a.storage.GetQuestByID(ctx, questID)
	a.require.NoError(err, "Failed to get quest from database")

	a.assert.False(quest.CreatedAt.IsZero(), "CreatedAt should not be zero")
	a.assert.False(quest.UpdatedAt.IsZero(), "UpdatedAt should not be zero")
	a.assert.True(quest.UpdatedAt.After(quest.CreatedAt) || quest.UpdatedAt.Equal(quest.CreatedAt),
		"UpdatedAt should be after or equal to CreatedAt")
}

// EventExists проверяет что событие существует
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

// EventCountEquals проверяет количество событий для агрегата
func (a *QuestAssertions) EventCountEquals(ctx context.Context, aggregateID uuid.UUID, expectedCount int64) {
	count, err := a.eventStorage.CountEventsByAggregateID(ctx, aggregateID)
	a.require.NoError(err, "Failed to count events")
	a.assert.Equal(expectedCount, count, "Event count should match expected")
}

// EventCountByTypeEquals проверяет количество событий по типу
func (a *QuestAssertions) EventCountByTypeEquals(ctx context.Context, eventType string, expectedCount int64) {
	count, err := a.eventStorage.CountEventsByType(ctx, eventType)
	a.require.NoError(err, "Failed to count events by type")
	a.assert.Equal(expectedCount, count, "Event count by type should match expected")
}

// WaitForEvents ждет появления событий (для асинхронных операций)
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

	// Финальная проверка
	count, err := a.eventStorage.CountEventsByAggregateID(ctx, aggregateID)
	a.require.NoError(err, "Failed to count events")
	a.assert.GreaterOrEqual(count, expectedCount, "Should have at least %d events after timeout", expectedCount)
}
