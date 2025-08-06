package contracts

import (
	"context"
	"testing"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// EventPublisherContractSuite defines contract tests that all EventPublisher implementations must pass
type EventPublisherContractSuite struct {
	suite.Suite
	publisher ports.EventPublisher
	ctx       context.Context
}

func (s *EventPublisherContractSuite) SetupTest() {
	// Clear state for MockEventPublisher before each test
	if mockPublisher, ok := s.publisher.(*MockEventPublisher); ok {
		mockPublisher.PublishedEvents = nil
		mockPublisher.PublishAsyncEvents = nil
		mockPublisher.PublishError = nil
	}
}

// MockEventPublisher for testing contract behavior
type MockEventPublisher struct {
	PublishedEvents    []ddd.DomainEvent
	PublishError       error
	PublishAsyncEvents []ddd.DomainEvent
}

func (m *MockEventPublisher) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	_ = ctx // unused in mock
	if m.PublishError != nil {
		return m.PublishError
	}
	m.PublishedEvents = append(m.PublishedEvents, events...)
	return nil
}

func (m *MockEventPublisher) PublishAsync(ctx context.Context, events ...ddd.DomainEvent) {
	_ = ctx // unused in mock
	m.PublishAsyncEvents = append(m.PublishAsyncEvents, events...)
}

// TestNullEventPublisherContract tests the NullEventPublisher implementation
func TestNullEventPublisherContract(t *testing.T) {
	s := &EventPublisherContractSuite{
		publisher: &ports.NullEventPublisher{},
		ctx:       context.Background(),
	}
	suite.Run(t, s)
}

func TestMockEventPublisherContract(t *testing.T) {
	s := &EventPublisherContractSuite{
		publisher: &MockEventPublisher{},
		ctx:       context.Background(),
	}
	suite.Run(t, s)
}

// EventPublisher contract tests

func (s *EventPublisherContractSuite) TestPublishSingleEvent() {
	// Create a test event
	event := quest.NewQuestCreated(
		uuid.New(),
		"test-creator",
	)

	// Contract: Publish should handle a single event without error
	err := s.publisher.Publish(s.ctx, event)
	s.Assert().NoError(err, "Publish should succeed with a single event")

	// For MockEventPublisher, verify the event was captured
	if mockPublisher, ok := s.publisher.(*MockEventPublisher); ok {
		s.Assert().Len(mockPublisher.PublishedEvents, 1, "Should have published exactly one event")
		s.Assert().Equal(event.GetID(), mockPublisher.PublishedEvents[0].GetID(), "Published event should match")
		s.Assert().Equal(event.GetName(), mockPublisher.PublishedEvents[0].GetName(), "Published event name should match")
	}
}

func (s *EventPublisherContractSuite) TestPublishMultipleEvents() {
	// Create multiple test events
	event1 := quest.NewQuestCreated(
		uuid.New(),
		"test-creator-1",
	)

	event2 := quest.NewQuestAssigned(
		uuid.New(),
		"test-assignee",
	)

	event3 := quest.NewQuestStatusChanged(
		uuid.New(),
		quest.StatusCreated,
		quest.StatusAssigned,
	)

	// Contract: Publish should handle multiple events without error
	err := s.publisher.Publish(s.ctx, event1, event2, event3)
	s.Assert().NoError(err, "Publish should succeed with multiple events")

	// For MockEventPublisher, verify all events were captured
	if mockPublisher, ok := s.publisher.(*MockEventPublisher); ok {
		s.Assert().Len(mockPublisher.PublishedEvents, 3, "Should have published exactly three events")

		// Verify each event
		publishedIDs := make(map[uuid.UUID]bool)
		for _, publishedEvent := range mockPublisher.PublishedEvents {
			publishedIDs[publishedEvent.GetID()] = true
		}

		s.Assert().True(publishedIDs[event1.GetID()], "Should have published first event")
		s.Assert().True(publishedIDs[event2.GetID()], "Should have published second event")
		s.Assert().True(publishedIDs[event3.GetID()], "Should have published third event")
	}
}

func (s *EventPublisherContractSuite) TestPublishNoEvents() {
	// Contract: Publish should handle empty event list without error
	err := s.publisher.Publish(s.ctx)
	s.Assert().NoError(err, "Publish should succeed with no events")

	// For MockEventPublisher, verify no events were added
	if mockPublisher, ok := s.publisher.(*MockEventPublisher); ok {
		// Length should remain the same as before (could be 0 or more from previous tests)
		initialLength := len(mockPublisher.PublishedEvents)

		err = s.publisher.Publish(s.ctx)
		s.Assert().NoError(err)
		s.Assert().Len(mockPublisher.PublishedEvents, initialLength, "Should not add any events when publishing empty list")
	}
}

func (s *EventPublisherContractSuite) TestPublishAsyncSingleEvent() {
	// Create a test event
	event := quest.NewQuestCreated(
		uuid.New(),
		"async-creator",
	)

	// Contract: PublishAsync should handle a single event without blocking or returning error
	// Since PublishAsync is asynchronous, we can't directly assert on errors, but it shouldn't panic
	s.Assert().NotPanics(func() {
		s.publisher.PublishAsync(s.ctx, event)
	}, "PublishAsync should not panic with a single event")

	// For MockEventPublisher, verify the event was captured
	if mockPublisher, ok := s.publisher.(*MockEventPublisher); ok {
		s.Assert().Len(mockPublisher.PublishAsyncEvents, 1, "Should have async published exactly one event")
		s.Assert().Equal(event.GetID(), mockPublisher.PublishAsyncEvents[0].GetID(), "Async published event should match")
	}
}

func (s *EventPublisherContractSuite) TestPublishAsyncMultipleEvents() {
	// Create multiple test events
	event1 := quest.NewQuestStatusChanged(
		uuid.New(),
		quest.StatusAssigned,
		quest.StatusInProgress,
	)

	event2 := quest.NewQuestStatusChanged(
		uuid.New(),
		quest.StatusInProgress,
		quest.StatusCompleted,
	)

	// Contract: PublishAsync should handle multiple events without blocking or panicking
	s.Assert().NotPanics(func() {
		s.publisher.PublishAsync(s.ctx, event1, event2)
	}, "PublishAsync should not panic with multiple events")

	// For MockEventPublisher, verify events were captured
	if mockPublisher, ok := s.publisher.(*MockEventPublisher); ok {
		// Should have at least these 2 events (might have more from previous tests)
		s.Assert().True(len(mockPublisher.PublishAsyncEvents) >= 2, "Should have async published at least two events")

		// Find our events in the list
		found1, found2 := false, false
		for _, publishedEvent := range mockPublisher.PublishAsyncEvents {
			if publishedEvent.GetID() == event1.GetID() {
				found1 = true
			}
			if publishedEvent.GetID() == event2.GetID() {
				found2 = true
			}
		}
		s.Assert().True(found1, "Should have async published first event")
		s.Assert().True(found2, "Should have async published second event")
	}
}

func (s *EventPublisherContractSuite) TestPublishAsyncNoEvents() {
	// Contract: PublishAsync should handle empty event list without panicking
	s.Assert().NotPanics(func() {
		s.publisher.PublishAsync(s.ctx)
	}, "PublishAsync should not panic with no events")
}

// Test that PublishAsync doesn't interfere with synchronous Publish
func (s *EventPublisherContractSuite) TestPublishSyncAndAsyncIndependence() {
	if mockPublisher, ok := s.publisher.(*MockEventPublisher); ok {
		// Clear any previous events
		mockPublisher.PublishedEvents = nil
		mockPublisher.PublishAsyncEvents = nil

		// Create different events for sync and async
		syncEvent := quest.NewQuestCreated(
			uuid.New(),
			"sync-creator",
		)

		asyncEvent := quest.NewQuestAssigned(
			uuid.New(),
			"async-assignee",
		)

		// Publish synchronously
		err := s.publisher.Publish(s.ctx, syncEvent)
		s.Assert().NoError(err)

		// Publish asynchronously
		s.publisher.PublishAsync(s.ctx, asyncEvent)

		// Contract: Sync and async events should be tracked separately
		s.Assert().Len(mockPublisher.PublishedEvents, 1, "Should have one sync published event")
		s.Assert().Len(mockPublisher.PublishAsyncEvents, 1, "Should have one async published event")

		s.Assert().Equal(syncEvent.GetID(), mockPublisher.PublishedEvents[0].GetID(), "Sync event should match")
		s.Assert().Equal(asyncEvent.GetID(), mockPublisher.PublishAsyncEvents[0].GetID(), "Async event should match")
	}
}

// Test error handling in Publish method
func (s *EventPublisherContractSuite) TestPublishErrorHandling() {
	if mockPublisher, ok := s.publisher.(*MockEventPublisher); ok {
		// Setup mock to return an error
		originalError := mockPublisher.PublishError
		mockPublisher.PublishError = &MockPublishError{message: "simulated publish failure"}
		defer func() {
			mockPublisher.PublishError = originalError
		}()

		event := quest.NewQuestCreated(
			uuid.New(),
			"error-creator",
		)

		// Contract: Publish should return the error when publisher fails
		err := s.publisher.Publish(s.ctx, event)
		s.Assert().Error(err, "Publish should return error when publisher fails")
		s.Assert().Contains(err.Error(), "simulated publish failure", "Error should contain expected message")
	}
}

// MockPublishError for testing error handling
type MockPublishError struct {
	message string
}

func (e *MockPublishError) Error() string {
	return e.message
}

// Test context handling
func (s *EventPublisherContractSuite) TestContextHandling() {
	// Create a context with timeout or cancellation
	ctx, cancel := context.WithCancel(s.ctx)

	event := quest.NewQuestCreated(
		uuid.New(),
		"context-creator",
	)

	// Contract: Publisher should accept and handle context properly
	err := s.publisher.Publish(ctx, event)
	s.Assert().NoError(err, "Publish should succeed with custom context")

	// Cancel context and test async publish
	cancel()

	// Contract: PublishAsync should handle cancelled context gracefully (not panic)
	s.Assert().NotPanics(func() {
		s.publisher.PublishAsync(ctx, event)
	}, "PublishAsync should not panic with cancelled context")
}
