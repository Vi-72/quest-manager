package quest_e2e_tests

import (
	"context"
	"time"

	"quest-manager/internal/generated/servers"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"

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

	// 3. Verify locations were stored in the database
	e2eAssertions := assertions.NewQuestE2EAssertions(s.Assert(), s.TestDIContainer.EventStorage)
	e2eAssertions.VerifyLocationsInDatabase(ctx, questRequest, s.TestDIContainer.LocationRepository)

	// 4. Verify events were processed (quest created successfully)
	e2eAssertions.VerifyQuestBasicProperties(dbQuest)

	// 5. Verify events were saved in the events table
	e2eAssertions.VerifyQuestEventsCreated(ctx, dbQuest.ID())

	// Verify location events were also created
	e2eAssertions.VerifyLocationEventsCreated(ctx, dbQuest)
}

// Test 2: Try to create quest via API with invalid coordinates, verify empty tables
func (s *E2ESuite) TestCreateQuestThroughAPIInvalidCoordinates() {
	ctx := context.Background()

	// Count initial records
	initialQuests, initialLocations, err := casesteps.CountInitialDatabaseRecords(ctx, s.TestDIContainer.QuestRepository, s.TestDIContainer.LocationRepository)
	s.Require().NoError(err)
	initialQuestCount := len(initialQuests)
	initialLocationCount := len(initialLocations)

	// Prepare request with invalid coordinates (latitude > 90)
	questRequest := testdatagenerators.InvalidCoordinatesQuestData()

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
	e2eAssertions := assertions.NewQuestE2EAssertions(s.Assert(), s.TestDIContainer.EventStorage)
	e2eAssertions.VerifyNoEventsCreated(ctx)

	// This is also implicitly verified by the quest count remaining the same
	s.Assert().Equal(initialQuestCount, len(finalQuests), "No new quest events should have been processed")
}
