package assertions

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	teststorage "quest-manager/tests/integration/core/storage"
)

type QuestE2EAssertions struct {
	assert       *assert.Assertions
	eventStorage *teststorage.EventStorage
}

func NewQuestE2EAssertions(a *assert.Assertions, eventStorage *teststorage.EventStorage) *QuestE2EAssertions {
	return &QuestE2EAssertions{
		assert:       a,
		eventStorage: eventStorage,
	}
}

// VerifyQuestEventsCreated checks that quest.created event exists for the quest
func (a *QuestE2EAssertions) VerifyQuestEventsCreated(ctx context.Context, questID uuid.UUID) {
	questEvents, err := a.eventStorage.GetEventsByAggregateID(ctx, questID)
	a.assert.NoError(err, "Should retrieve events for the created quest")
	a.assert.GreaterOrEqual(len(questEvents), 1, "Should have at least one event for quest creation")

	// Check that quest.created event exists
	var questCreatedEventFound bool
	for _, event := range questEvents {
		if event.EventType == "quest.created" {
			questCreatedEventFound = true
			a.assert.Equal(questID.String(), event.AggregateID, "Event aggregate ID should match quest ID")
			a.assert.NotEmpty(event.Data, "Event data should not be empty")
			a.assert.True(event.CreatedAt.Before(time.Now().Add(time.Second)), "Event creation time should be recent")
			break
		}
	}
	a.assert.True(questCreatedEventFound, "Should find quest.created event in events table")
}

// VerifyQuestAssignmentEvents checks that quest.assigned and quest.status_changed events exist
func (a *QuestE2EAssertions) VerifyQuestAssignmentEvents(ctx context.Context, questID uuid.UUID, createdAt time.Time) {
	questEvents, err := a.eventStorage.GetEventsByAggregateID(ctx, questID)
	a.assert.NoError(err, "Should retrieve events for the assigned quest")
	a.assert.GreaterOrEqual(len(questEvents), 2, "Should have at least 2 events (created + assigned)")

	// Check that quest.assigned and quest.status_changed events exist
	var questAssignedEventFound bool
	var questStatusChangedEventFound bool
	for _, event := range questEvents {
		if event.EventType == "quest.assigned" {
			questAssignedEventFound = true
			a.assert.Equal(questID.String(), event.AggregateID, "Assign event aggregate ID should match quest ID")
			a.assert.NotEmpty(event.Data, "Assign event data should not be empty")
			a.assert.True(event.CreatedAt.After(createdAt), "Assign event should be after quest creation")
		}
		if event.EventType == "quest.status_changed" {
			questStatusChangedEventFound = true
			a.assert.Equal(questID.String(), event.AggregateID, "Status change event aggregate ID should match quest ID")
			a.assert.NotEmpty(event.Data, "Status change event data should not be empty")
			a.assert.True(event.CreatedAt.After(createdAt), "Status change event should be after quest creation")
		}
	}
	a.assert.True(questAssignedEventFound, "Should find quest.assigned event in events table")
	a.assert.True(questStatusChangedEventFound, "Should find quest.status_changed event in events table")
}

// VerifyLocationEventsCreated checks that location events exist for quest locations
func (a *QuestE2EAssertions) VerifyLocationEventsCreated(ctx context.Context, quest quest.Quest) {
	// Verify location events were also created (if location IDs are set)
	if quest.TargetLocationID != nil {
		targetLocationEvents, err := a.eventStorage.GetEventsByAggregateID(ctx, *quest.TargetLocationID)
		a.assert.NoError(err, "Should retrieve events for target location")
		a.assert.GreaterOrEqual(len(targetLocationEvents), 1, "Should have at least one event for target location")
	}

	if quest.ExecutionLocationID != nil {
		executionLocationEvents, err := a.eventStorage.GetEventsByAggregateID(ctx, *quest.ExecutionLocationID)
		a.assert.NoError(err, "Should retrieve events for execution location")
		a.assert.GreaterOrEqual(len(executionLocationEvents), 1, "Should have at least one event for execution location")
	}
}

// VerifyNoEventsCreated checks that no events were created during this test
func (a *QuestE2EAssertions) VerifyNoEventsCreated(ctx context.Context) {
	allEvents, err := a.eventStorage.GetAllEvents(ctx)
	a.assert.NoError(err, "Should retrieve all events from database")

	// Filter events by creation time to find only events created during this test
	recentEvents := 0
	testStartTime := time.Now().Add(-5 * time.Second) // Events created in last 5 seconds
	for _, event := range allEvents {
		if event.CreatedAt.After(testStartTime) {
			recentEvents++
		}
	}
	a.assert.Equal(0, recentEvents, "No new events should have been created for invalid quest creation")
}

// VerifyLocationsInDatabase checks that quest locations are properly stored in database
func (a *QuestE2EAssertions) VerifyLocationsInDatabase(ctx context.Context, questRequest v1.CreateQuestRequest, locationRepo interface {
	FindAll(ctx context.Context) ([]*location.Location, error)
}) {
	allLocations, err := locationRepo.FindAll(ctx)
	a.assert.NoError(err, "Should retrieve all locations from database")
	a.assert.GreaterOrEqual(len(allLocations), 2, "Should have at least 2 locations (target and execution)")

	var targetLocationFound, executionLocationFound bool
	for _, loc := range allLocations {
		if a.coordinatesMatch(loc.Coordinate.Lat, loc.Coordinate.Lon,
			float64(questRequest.TargetLocation.Latitude), float64(questRequest.TargetLocation.Longitude)) {
			targetLocationFound = true
		}
		if a.coordinatesMatch(loc.Coordinate.Lat, loc.Coordinate.Lon,
			float64(questRequest.ExecutionLocation.Latitude), float64(questRequest.ExecutionLocation.Longitude)) {
			executionLocationFound = true
		}
	}
	a.assert.True(targetLocationFound, "Target location should be found in database")
	a.assert.True(executionLocationFound, "Execution location should be found in database")
}

// coordinatesMatch checks if two coordinate pairs match within a small tolerance
func (a *QuestE2EAssertions) coordinatesMatch(lat1, lon1, lat2, lon2 float64) bool {
	return math.Abs(lat1-lat2) < 0.0001 && math.Abs(lon1-lon2) < 0.0001
}

// VerifyQuestBasicProperties checks basic quest properties after creation
func (a *QuestE2EAssertions) VerifyQuestBasicProperties(quest quest.Quest) {
	a.assert.Equal("created", string(quest.Status), "Quest should be in created status")
	a.assert.NotNil(quest.CreatedAt, "Quest should have creation timestamp")
	a.assert.NotEqual(quest.ID().String(), "", "Quest should have valid ID (indicates event processing)")
	a.assert.True(quest.CreatedAt.Before(time.Now().Add(time.Second)), "Creation time should be recent")
	a.assert.True(quest.UpdatedAt.Equal(quest.CreatedAt) || quest.UpdatedAt.After(quest.CreatedAt),
		"Updated time should be equal or after created time for new quest")
}
