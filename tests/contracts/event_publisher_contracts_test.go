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
		mockPublisher.PublishError = nil
	}
}

// MockEventPublisher for testing contract behavior
type MockEventPublisher struct {
	PublishedEvents []ddd.DomainEvent
	PublishError    error
}

func (m *MockEventPublisher) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	_ = ctx // unused in mock
	if m.PublishError != nil {
		return m.PublishError
	}
	m.PublishedEvents = append(m.PublishedEvents, events...)
	return nil
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
		uuid.New(), // test-assignee as UUID
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

	cancel()
}
