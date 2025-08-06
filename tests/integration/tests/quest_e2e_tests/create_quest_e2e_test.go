package quest_e2e_tests

import (
	"context"
	"math"
	"time"

	"quest-manager/internal/generated/servers"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"

	"github.com/google/uuid"
)

// Test 1: Create quest via API with locations, verify database and events
func (s *E2ESuite) TestCreateQuestThroughAPISuccess() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Prepare request data
	questRequest := servers.CreateQuestRequest{
		Title:           "E2E Test Quest",
		Description:     "End-to-end test quest creation",
		Difficulty:      servers.CreateQuestRequestDifficultyMedium,
		Reward:          3,
		DurationMinutes: 60,
		TargetLocation: servers.Coordinate{
			Latitude:  55.7558,
			Longitude: 37.6176,
		},
		ExecutionLocation: servers.Coordinate{
			Latitude:  55.7520,
			Longitude: 37.6175,
		},
		Equipment: &[]string{"passport", "camera"},
		Skills:    &[]string{"sightseeing", "photography"},
	}

	// Execute HTTP request
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert HTTP response
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)
	httpAssertions.QuestArraysNotNull(createdQuest)

	// 2. Verify quest exists in database
	dbQuest, err := s.TestDIContainer.QuestRepository.GetByID(ctx, uuid.MustParse(createdQuest.Id))
	s.Require().NoError(err, "Quest should exist in database")
	s.Assert().Equal(questRequest.Title, dbQuest.Title, "Quest title should match")
	s.Assert().Equal(questRequest.Description, dbQuest.Description, "Quest description should match")
	s.Assert().Equal(string(questRequest.Difficulty), string(dbQuest.Difficulty), "Quest difficulty should match")

	// 3. Verify locations are in the database
	s.Assert().InDelta(questRequest.TargetLocation.Latitude, dbQuest.TargetLocation.Lat, 0.0001, "Target location latitude should match")
	s.Assert().InDelta(questRequest.TargetLocation.Longitude, dbQuest.TargetLocation.Lon, 0.0001, "Target location longitude should match")
	s.Assert().InDelta(questRequest.ExecutionLocation.Latitude, dbQuest.ExecutionLocation.Lat, 0.0001, "Execution location latitude should match")
	s.Assert().InDelta(questRequest.ExecutionLocation.Longitude, dbQuest.ExecutionLocation.Lon, 0.0001, "Execution location longitude should match")

	allLocations, err := s.TestDIContainer.LocationRepository.FindAll(ctx)
	s.Require().NoError(err, "Should retrieve all locations from database")
	s.Assert().GreaterOrEqual(len(allLocations), 2, "Should have at least 2 locations (target and execution)")

	var targetLocationFound, executionLocationFound bool
	for _, loc := range allLocations {
		if math.Abs(loc.Coordinate.Lat-float64(questRequest.TargetLocation.Latitude)) < 0.0001 &&
			math.Abs(loc.Coordinate.Lon-float64(questRequest.TargetLocation.Longitude)) < 0.0001 {
			targetLocationFound = true
		}
		if math.Abs(loc.Coordinate.Lat-float64(questRequest.ExecutionLocation.Latitude)) < 0.0001 &&
			math.Abs(loc.Coordinate.Lon-float64(questRequest.ExecutionLocation.Longitude)) < 0.0001 {
			executionLocationFound = true
		}
	}
	s.Assert().True(targetLocationFound, "Target location should be found in database")
	s.Assert().True(executionLocationFound, "Execution location should be found in database")

	// 4. Verify events were processed (quest created successfully)
	s.Assert().Equal("created", string(dbQuest.Status), "Quest should be in created status")
	s.Assert().NotNil(dbQuest.CreatedAt, "Quest should have creation timestamp")
	s.Assert().NotEqual(dbQuest.ID().String(), "", "Quest should have valid ID (indicates event processing)")
	s.Assert().True(dbQuest.CreatedAt.Before(time.Now().Add(time.Second)), "Creation time should be recent")
	s.Assert().True(dbQuest.UpdatedAt.Equal(dbQuest.CreatedAt) || dbQuest.UpdatedAt.After(dbQuest.CreatedAt), "Updated time should be equal or after created time for new quest")

	// 5. Verify events were saved in the events table
	questEvents, err := s.TestDIContainer.EventStorage.GetEventsByAggregateID(ctx, dbQuest.ID())
	s.Require().NoError(err, "Should retrieve events for the created quest")
	s.Assert().GreaterOrEqual(len(questEvents), 1, "Should have at least one event for quest creation")

	// Check that quest.created event exists
	var questCreatedEventFound bool
	for _, event := range questEvents {
		if event.EventType == "quest.created" {
			questCreatedEventFound = true
			s.Assert().Equal(dbQuest.ID().String(), event.AggregateID, "Event aggregate ID should match quest ID")
			s.Assert().NotEmpty(event.Data, "Event data should not be empty")
			s.Assert().True(event.CreatedAt.Before(time.Now().Add(time.Second)), "Event creation time should be recent")
			break
		}
	}
	s.Assert().True(questCreatedEventFound, "Should find quest.created event in events table")

	// Verify location events were also created (if location IDs are set)
	if dbQuest.TargetLocationID != nil {
		targetLocationEvents, err := s.TestDIContainer.EventStorage.GetEventsByAggregateID(ctx, *dbQuest.TargetLocationID)
		s.Require().NoError(err, "Should retrieve events for target location")
		s.Assert().GreaterOrEqual(len(targetLocationEvents), 1, "Should have at least one event for target location")
	}

	if dbQuest.ExecutionLocationID != nil {
		executionLocationEvents, err := s.TestDIContainer.EventStorage.GetEventsByAggregateID(ctx, *dbQuest.ExecutionLocationID)
		s.Require().NoError(err, "Should retrieve events for execution location")
		s.Assert().GreaterOrEqual(len(executionLocationEvents), 1, "Should have at least one event for execution location")
	}
}

// Test 2: Try to create quest via API with invalid coordinates, verify empty tables
func (s *E2ESuite) TestCreateQuestThroughAPIInvalidCoordinates() {
	ctx := context.Background()

	// Count initial records
	initialQuests, err := s.TestDIContainer.QuestRepository.FindAll(ctx)
	s.Require().NoError(err)
	initialQuestCount := len(initialQuests)

	initialLocations, err := s.TestDIContainer.LocationRepository.FindAll(ctx)
	s.Require().NoError(err)
	initialLocationCount := len(initialLocations)

	// Prepare request with invalid coordinates (latitude > 90)
	questRequest := servers.CreateQuestRequest{
		Title:           "Invalid Quest",
		Description:     "Quest with invalid coordinates",
		Difficulty:      servers.CreateQuestRequestDifficultyEasy,
		Reward:          1,
		DurationMinutes: 30,
		TargetLocation: servers.Coordinate{
			Latitude:  95.0, // Invalid: > 90
			Longitude: 37.6176,
		},
		ExecutionLocation: servers.Coordinate{
			Latitude:  55.7520,
			Longitude: 37.6175,
		},
	}

	// Execute HTTP request
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert HTTP response shows validation error
	s.Assert().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(400, createResp.StatusCode, "Should return bad request for invalid coordinates")

	// Wait for any async processing
	time.Sleep(100 * time.Millisecond)

	// 2. Verify quest table has no new records
	finalQuests, err := s.TestDIContainer.QuestRepository.FindAll(ctx)
	s.Require().NoError(err)
	s.Assert().Equal(initialQuestCount, len(finalQuests), "Quest count should remain unchanged")

	// 3. Verify location table has no new records
	finalLocations, err := s.TestDIContainer.LocationRepository.FindAll(ctx)
	s.Require().NoError(err)
	s.Assert().Equal(initialLocationCount, len(finalLocations), "Location count should remain unchanged")

	// 4. Verify no new events were created in events table
	allEvents, err := s.TestDIContainer.EventStorage.GetAllEvents(ctx)
	s.Require().NoError(err, "Should retrieve all events from database")

	// Filter events by creation time to find only events created during this test
	recentEvents := 0
	testStartTime := time.Now().Add(-5 * time.Second) // Events created in last 5 seconds
	for _, event := range allEvents {
		if event.CreatedAt.After(testStartTime) {
			recentEvents++
		}
	}
	s.Assert().Equal(0, recentEvents, "No new events should have been created for invalid quest creation")

	// This is also implicitly verified by the quest count remaining the same
	s.Assert().Equal(initialQuestCount, len(finalQuests), "No new quest events should have been processed")
}
