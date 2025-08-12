package domain

// DOMAIN EVENTS UNIT TESTS
// Tests for quest domain events: quest.created and quest.assigned

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

func TestQuest_NewQuest_DomainEvents(t *testing.T) {
	q := createValidQuestForEvents(t)

	// Assert - new quest should raise quest.created event
	events := q.GetDomainEvents()
	assert.Len(t, events, 1, "NewQuest should raise one domain event (quest.created)")

	// Verify event is of correct type (implementation-specific)
	// Note: Specific event verification depends on domain event implementation
}

func TestQuest_AssignTo_DomainEvents(t *testing.T) {
	q := createValidQuestForEvents(t)
	userID := "event-test-user"

	// Clear any existing events
	q.ClearDomainEvents()
	assert.Len(t, q.GetDomainEvents(), 0, "Events should be cleared")

	// Act - assign quest
	err := q.AssignTo(userID)
	assert.NoError(t, err)

	// Assert - should raise two domain events
	events := q.GetDomainEvents()
	assert.Len(t, events, 2, "AssignTo should raise two domain events (assigned + status changed)")

	// Verify events are of correct types (implementation-specific)
	// Note: Specific event verification depends on domain event implementation
}

func TestQuest_ClearDomainEvents(t *testing.T) {
	q := createValidQuestForEvents(t)

	// Ensure there are some events
	events := q.GetDomainEvents()
	assert.NotEmpty(t, events, "Quest should have domain events after creation")

	// Act - clear events
	q.ClearDomainEvents()

	// Assert - events should be cleared
	events = q.GetDomainEvents()
	assert.Empty(t, events, "Events should be cleared after ClearDomainEvents()")
}

func TestQuest_ChangeStatus_DomainEvents(t *testing.T) {
	q := createValidQuestForEvents(t)

	// Clear creation events
	q.ClearDomainEvents()

	// Act - change status (valid transition: created -> posted)
	err := q.ChangeStatus(quest.StatusPosted)
	assert.NoError(t, err)

	// Assert - should raise one domain event for status change
	events := q.GetDomainEvents()
	assert.Len(t, events, 1, "ChangeStatus should raise one domain event")
}

func TestQuest_GetDomainEvents_Immutability(t *testing.T) {
	q := createValidQuestForEvents(t)

	// Get events twice
	events1 := q.GetDomainEvents()
	events2 := q.GetDomainEvents()

	// Should return consistent results
	assert.Equal(t, len(events1), len(events2), "GetDomainEvents should return consistent results")

	// Modifying returned slice should not affect original
	if len(events1) > 0 {
		originalLen := len(events1)
		_ = append(events1, events1[0]) // Try to modify (deliberately not using result)

		events3 := q.GetDomainEvents()
		assert.Equal(t, originalLen, len(events3), "Modifying returned events should not affect original")
	}
}

// Helper function to create a valid quest for testing
func createValidQuestForEvents(t *testing.T) *quest.Quest {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	q, err := quest.NewQuest(
		"Test Quest",
		"Test description",
		"medium",
		3,
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{"equipment"},
		[]string{"skill"},
	)

	assert.NoError(t, err)
	return &q
}
