package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
)

type EventRepositoryTestSuite struct {
	RepositoryTestSuite
	eventRepo ports.EventPublisher
}

func (suite *EventRepositoryTestSuite) SetupSuite() {
	// Call parent setup
	suite.RepositoryTestSuite.SetupSuite()

	// Create event repository
	eventRepo, err := eventrepo.NewRepository(suite.tracker, 5) // 5 goroutines limit
	suite.Require().NoError(err)
	suite.eventRepo = eventRepo
}

func (suite *EventRepositoryTestSuite) TestPublish_Success() {
	ctx := context.Background()

	// Arrange - create test events
	questID := uuid.New()
	events := []ddd.DomainEvent{
		quest.NewQuestCreated(questID, "test-creator"),
		quest.NewQuestStatusChanged(questID, quest.StatusCreated, quest.StatusPosted),
	}

	// Act - publish events synchronously
	err := suite.eventRepo.Publish(ctx, events...)

	// Assert
	suite.Require().NoError(err)

	// Verify events were saved to database
	suite.verifyEventsInDatabase(2)
}

func (suite *EventRepositoryTestSuite) TestPublish_EmptyEvents() {
	ctx := context.Background()

	// Act - publish empty events
	err := suite.eventRepo.Publish(ctx)

	// Assert
	suite.Require().NoError(err)

	// Verify no events in database
	suite.verifyEventsInDatabase(0)
}

func (suite *EventRepositoryTestSuite) TestPublish_SingleEvent() {
	ctx := context.Background()

	// Arrange - create single event
	questID := uuid.New()
	event := quest.NewQuestCreated(questID, "test-creator")

	// Act - publish single event
	err := suite.eventRepo.Publish(ctx, event)

	// Assert
	suite.Require().NoError(err)

	// Verify event was saved
	suite.verifyEventsInDatabase(1)
}

func (suite *EventRepositoryTestSuite) TestPublish_MultipleEventTypes() {
	ctx := context.Background()

	// Arrange - create different types of events
	questID := uuid.New()

	questCreated := quest.NewQuestCreated(questID, "test-creator")
	questAssigned := quest.NewQuestAssigned(questID, "test-user")
	statusChanged := quest.NewQuestStatusChanged(questID, quest.StatusCreated, quest.StatusAssigned)

	// Act - publish multiple different events
	err := suite.eventRepo.Publish(ctx, questCreated, questAssigned, statusChanged)

	// Assert
	suite.Require().NoError(err)

	// Verify all events were saved
	suite.verifyEventsInDatabase(3)
}

func (suite *EventRepositoryTestSuite) TestPublishAsync_Success() {
	ctx := context.Background()

	// Arrange - create test events
	questID := uuid.New()
	events := []ddd.DomainEvent{
		quest.NewQuestCreated(questID, "test-creator"),
		quest.NewQuestStatusChanged(questID, quest.StatusCreated, quest.StatusPosted),
	}

	// Act - publish events asynchronously
	suite.eventRepo.PublishAsync(ctx, events...)

	// Wait for async processing to complete
	time.Sleep(100 * time.Millisecond)

	// Assert - verify events were saved to database
	suite.verifyEventsInDatabase(2)
}

func (suite *EventRepositoryTestSuite) TestPublishAsync_EmptyEvents() {
	ctx := context.Background()

	// Act - publish empty events async
	suite.eventRepo.PublishAsync(ctx)

	// Wait a bit (though nothing should happen)
	time.Sleep(50 * time.Millisecond)

	// Assert - verify no events in database
	suite.verifyEventsInDatabase(0)
}

func (suite *EventRepositoryTestSuite) TestPublishAsync_HighVolume() {
	ctx := context.Background()

	// Arrange - create many events to test goroutine limiting
	var events []ddd.DomainEvent
	for i := 0; i < 20; i++ {
		questID := uuid.New()
		events = append(events, quest.NewQuestCreated(questID, "test-creator"))
	}

	// Act - publish many events async
	suite.eventRepo.PublishAsync(ctx, events...)

	// Wait for all async processing to complete
	time.Sleep(500 * time.Millisecond)

	// Assert - verify all events were saved
	suite.verifyEventsInDatabase(20)
}

func (suite *EventRepositoryTestSuite) TestPublish_WithTransaction() {
	ctx := context.Background()

	// Arrange - start transaction manually
	err := suite.tracker.Begin(ctx)
	suite.Require().NoError(err)

	questID := uuid.New()
	event := quest.NewQuestCreated(questID, "test-creator")

	// Act - publish event within transaction
	err = suite.eventRepo.Publish(ctx, event)
	suite.Require().NoError(err)

	// Before commit - verify event is not visible from outside transaction
	// (This would require a separate connection to verify, simplified here)

	// Commit transaction
	err = suite.tracker.Commit(ctx)
	suite.Require().NoError(err)

	// Assert - verify event is now visible
	suite.verifyEventsInDatabase(1)
}

func (suite *EventRepositoryTestSuite) TestPublish_TransactionRollback() {
	ctx := context.Background()

	// Arrange - start transaction
	err := suite.tracker.Begin(ctx)
	suite.Require().NoError(err)

	questID := uuid.New()
	event := quest.NewQuestCreated(questID, "test-creator")

	// Act - publish event within transaction
	err = suite.eventRepo.Publish(ctx, event)
	suite.Require().NoError(err)

	// Rollback transaction
	err = suite.tracker.Rollback()
	suite.Require().NoError(err)

	// Assert - verify event was not saved due to rollback
	suite.verifyEventsInDatabase(0)
}

func (suite *EventRepositoryTestSuite) TestPublish_ComplexDomainScenario() {
	ctx := context.Background()

	// Arrange - simulate complex domain scenario
	userID := "test-user-123"

	// Create quest
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
	suite.Require().NoError(err)

	// Get events from quest creation and publish them immediately
	creationEvents := q.GetDomainEvents()
	err = suite.eventRepo.Publish(ctx, creationEvents...)
	suite.Require().NoError(err)

	// Clear events after publishing
	q.ClearDomainEvents()

	// Assign quest (generates more events)
	err = q.AssignTo(userID)
	suite.Require().NoError(err)

	assignmentEvents := q.GetDomainEvents()

	// Act - publish assignment events
	err = suite.eventRepo.Publish(ctx, assignmentEvents...)

	// Assert
	suite.Require().NoError(err)
	totalExpectedEvents := len(creationEvents) + len(assignmentEvents)
	suite.verifyEventsInDatabase(totalExpectedEvents)
}

// Helper methods
func (suite *EventRepositoryTestSuite) verifyEventsInDatabase(expectedCount int) {
	var count int64
	err := suite.db.Table("events").Count(&count).Error
	suite.Require().NoError(err)
	suite.Equal(int64(expectedCount), count, "Expected %d events in database, found %d", expectedCount, count)
}

func (suite *EventRepositoryTestSuite) cleanEventDatabase() {
	suite.db.Exec("DELETE FROM events")
}

func (suite *EventRepositoryTestSuite) SetupTest() {
	// Call parent setup
	suite.RepositoryTestSuite.SetupTest()

	// Also clean events table
	suite.cleanEventDatabase()
}
