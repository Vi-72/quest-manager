//go:build integration

package repository

// REPOSITORY LAYER INTEGRATION TESTS
// Tests for repository implementations and database interactions

import (
	"context"
	"time"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/pkg/ddd"
	teststorage "quest-manager/tests/integration/core/storage"

	"github.com/google/uuid"
)

func (s *Suite) TestEventRepository_Publish_SingleEvent() {
	ctx := context.Background()

	// Pre-condition - create a domain event
	aggregateID := uuid.New()
	event := s.createTestEvent("test.event", aggregateID, map[string]interface{}{
		"message": "test event data",
		"value":   42,
	})

	// Act - publish event
	err := s.TestDIContainer.EventPublisher.Publish(ctx, event)

	// Assert
	s.Require().NoError(err)

	// Verify event was stored (using our event storage helper)
	eventStorage := teststorage.NewEventStorage(s.TestDIContainer.DB)
	events, err := eventStorage.GetEventsByType(ctx, "test.event")
	s.Require().NoError(err)
	s.Len(events, 1)

	storedEvent := events[0]
	s.Equal("test.event", storedEvent.EventType)
	s.Equal(aggregateID.String(), storedEvent.AggregateID)
	s.Contains(storedEvent.Data, "test event data")
	s.NotEmpty(storedEvent.ID)
	// Не проверяем конкретное ID события, так как оно генерируется случайно
}

func (s *Suite) TestEventRepository_Publish_MultipleEvents() {
	ctx := context.Background()

	// Pre-condition - create multiple domain events
	aggregateID := uuid.New()
	event1 := s.createTestEvent("quest.created", aggregateID, map[string]interface{}{
		"title": "Test Quest 1",
	})
	event2 := s.createTestEvent("quest.status_changed", aggregateID, map[string]interface{}{
		"from": "created",
		"to":   "posted",
	})
	event3 := s.createTestEvent("quest.assigned", aggregateID, map[string]interface{}{
		"assignee": "user123",
	})

	// Act - publish multiple events
	err := s.TestDIContainer.EventPublisher.Publish(ctx, event1, event2, event3)

	// Assert
	s.Require().NoError(err)

	// Verify all events were stored
	eventStorage := teststorage.NewEventStorage(s.TestDIContainer.DB)

	allEvents, err := eventStorage.GetEventsByAggregateID(ctx, aggregateID)
	s.Require().NoError(err)
	s.Len(allEvents, 3)

	// Verify events are in correct order (только если события действительно сохранились)
	if len(allEvents) >= 3 {
		s.Equal("quest.created", allEvents[0].EventType)
		s.Equal("quest.status_changed", allEvents[1].EventType)
		s.Equal("quest.assigned", allEvents[2].EventType)
	} else {
		s.Fail("Not enough events saved", "Expected 3 events, got %d", len(allEvents))
	}
}

func (s *Suite) TestEventRepository_Publish_EmptyEvents() {
	ctx := context.Background()

	// Act - publish empty events slice
	err := s.TestDIContainer.EventPublisher.Publish(ctx)

	// Assert - should not error
	s.Require().NoError(err)

	// Verify no events were stored
	eventStorage := teststorage.NewEventStorage(s.TestDIContainer.DB)
	count, err := eventStorage.CountEvents(ctx)
	s.Require().NoError(err)
	s.Equal(int64(0), count)
}

func (s *Suite) TestEventRepository_Publish_WithTransaction() {
	ctx := context.Background()

	// Pre-condition - create events
	aggregateID := uuid.New()
	event1 := s.createTestEvent("transaction.test", aggregateID, map[string]interface{}{
		"step": "1",
	})
	event2 := s.createTestEvent("transaction.test", aggregateID, map[string]interface{}{
		"step": "2",
	})

	// Act - publish events within transaction
	err := s.TestDIContainer.EventPublisher.Publish(ctx, event1, event2)

	// Assert
	s.Require().NoError(err)

	// Within transaction, events should be visible
	eventStorage := teststorage.NewEventStorage(s.TestDIContainer.DB)
	events, err := eventStorage.GetEventsByType(ctx, "transaction.test")
	s.Require().NoError(err)
	s.Len(events, 2)

	// Transaction will be rolled back in TearDownTest
	// In next test, events should not be visible (tested implicitly)
}

func (s *Suite) TestEventRepository_Publish_ComplexDomainScenario() {
	ctx := context.Background()

	// Pre-condition - create a real quest and trigger domain events
	q := s.createTestQuest("Complex Domain Test", "medium")

	// Act - perform domain operations that generate events
	// 1. Quest is created (triggers quest.created event)
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	// Publish domain events from quest creation
	err = s.TestDIContainer.EventPublisher.Publish(ctx, q.GetDomainEvents()...)
	s.Require().NoError(err)
	q.ClearDomainEvents()

	// 2. Change quest status (triggers quest.status_changed event)
	err = q.ChangeStatus(quest.StatusPosted)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)
	err = s.TestDIContainer.EventPublisher.Publish(ctx, q.GetDomainEvents()...)
	s.Require().NoError(err)
	q.ClearDomainEvents()

	// 3. Assign quest (triggers quest.assigned event)
	err = q.AssignTo("test-user")
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)
	err = s.TestDIContainer.EventPublisher.Publish(ctx, q.GetDomainEvents()...)
	s.Require().NoError(err)

	// Assert - verify complete event sequence
	eventStorage := teststorage.NewEventStorage(s.TestDIContainer.DB)
	questEvents, err := eventStorage.GetEventsByAggregateID(ctx, q.ID())
	s.Require().NoError(err)
	s.GreaterOrEqual(len(questEvents), 3) // At least 3 events

	// Verify event types
	eventTypes := make([]string, len(questEvents))
	for i, event := range questEvents {
		eventTypes[i] = event.EventType
	}

	// Should contain quest lifecycle events
	s.Contains(eventTypes, "quest.created")
	s.Contains(eventTypes, "quest.status_changed")
	s.Contains(eventTypes, "quest.assigned")
}

// ==========================================
// POSTGRESQL-SPECIFIC TESTS
// ==========================================

func (s *Suite) TestEventRepository_PostgreSQL_JSONEventData() {
	ctx := context.Background()

	// Test PostgreSQL JSON handling for complex event data
	aggregateID := uuid.New()
	complexData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":          "user123",
			"name":        "John Doe",
			"preferences": []string{"outdoor", "challenging"},
		},
		"quest": map[string]interface{}{
			"difficulty": "hard",
			"reward":     500,
			"location": map[string]interface{}{
				"lat": 55.7558,
				"lon": 37.6173,
			},
		},
		"metadata": map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"version":   "1.0",
		},
	}

	event := s.createTestEvent("complex.json.test", aggregateID, complexData)

	// Act - publish event with complex JSON
	err := s.TestDIContainer.EventPublisher.Publish(ctx, event)
	s.Require().NoError(err)

	// Assert - verify complex JSON is preserved
	eventStorage := teststorage.NewEventStorage(s.TestDIContainer.DB)
	events, err := eventStorage.GetEventsByType(ctx, "complex.json.test")
	s.Require().NoError(err)
	s.Len(events, 1)

	storedEvent := events[0]
	s.Contains(storedEvent.Data, "user123")
	s.Contains(storedEvent.Data, "John Doe")
	s.Contains(storedEvent.Data, "hard")
	s.Contains(storedEvent.Data, "55.7558")
}

func (s *Suite) TestEventRepository_PostgreSQL_ConcurrentEventPublishing() {
	ctx := context.Background()

	// Test PostgreSQL handling of concurrent event publishing
	aggregateID1 := uuid.New()
	aggregateID2 := uuid.New()

	event1 := s.createTestEvent("concurrent.test", aggregateID1, map[string]interface{}{
		"thread": "1",
	})
	event2 := s.createTestEvent("concurrent.test", aggregateID2, map[string]interface{}{
		"thread": "2",
	})

	// Act - publish events (simulating concurrent access)
	err1 := s.TestDIContainer.EventPublisher.Publish(ctx, event1)
	err2 := s.TestDIContainer.EventPublisher.Publish(ctx, event2)

	// Assert
	s.Require().NoError(err1)
	s.Require().NoError(err2)

	// Verify both events were stored correctly
	eventStorage := teststorage.NewEventStorage(s.TestDIContainer.DB)
	events, err := eventStorage.GetEventsByType(ctx, "concurrent.test")
	s.Require().NoError(err)
	s.Len(events, 2)

	// Verify data integrity
	aggregateIDs := make(map[string]bool)
	for _, event := range events {
		aggregateIDs[event.AggregateID] = true
	}

	// Более информативные проверки
	s.Contains(aggregateIDs, aggregateID1.String(), "Should find first aggregate ID")
	s.Contains(aggregateIDs, aggregateID2.String(), "Should find second aggregate ID")
	s.Len(aggregateIDs, 2, "Should have exactly 2 different aggregate IDs")
}

// ==========================================
// HELPER METHODS
// ==========================================

// TestEvent implements DomainEvent for testing
type TestEvent struct {
	ddd.BaseEvent
	Data map[string]interface{}
}

func (s *Suite) createTestEvent(eventType string, aggregateID uuid.UUID, data map[string]interface{}) ddd.DomainEvent {
	return TestEvent{
		BaseEvent: ddd.NewBaseEvent(aggregateID, eventType),
		Data:      data,
	}
}
