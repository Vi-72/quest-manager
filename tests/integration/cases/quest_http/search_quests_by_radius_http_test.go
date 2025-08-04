package quest_http

import (
	"context"
	"encoding/json"
	"net/http"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

func (s *Suite) TestSearchQuestsByRadiusHTTP() {
	ctx := context.Background()

	// Pre-condition - create quest at specific location via handler (for setup)
	centerLocation := testdatagenerators.DefaultTestCoordinate() // Moscow center: 55.7558, 37.6176
	nearLocation := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.01,  // ~1km away
		Lon: centerLocation.Longitude() + 0.01, // ~1km away
	}
	farLocation := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.1,  // ~10km away
		Lon: centerLocation.Longitude() + 0.1, // ~10km away
	}

	// Create quest near the center (within 5km radius)
	nearQuestData := testdatagenerators.SimpleQuestData(
		"Near Quest", "Quest near center", "easy", 1, 30, nearLocation, nearLocation)
	nearQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, nearQuestData)
	s.Require().NoError(err)

	// Create quest far from center (outside 5km radius)
	farQuestData := testdatagenerators.SimpleQuestData(
		"Far Quest", "Quest far from center", "medium", 2, 60, farLocation, farLocation)
	farQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, farQuestData)
	s.Require().NoError(err)

	// Act - search for quests within 5km radius via HTTP API
	radiusKm := 5.0
	searchReq := casesteps.SearchQuestsByRadiusHTTPRequest(
		centerLocation.Latitude(), centerLocation.Longitude(), radiusKm)
	searchResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, searchReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, searchResp.StatusCode)

	// Parse response
	var foundQuests []servers.Quest
	err = json.Unmarshal([]byte(searchResp.Body), &foundQuests)
	s.Require().NoError(err)

	// Should find the near quest but not the far quest
	nearQuestFound := false
	farQuestFound := false
	nearQuestID := nearQuest.ID().String()
	farQuestID := farQuest.ID().String()

	for _, q := range foundQuests {
		if q.Id == nearQuestID {
			nearQuestFound = true
		}
		if q.Id == farQuestID {
			farQuestFound = true
		}
	}

	s.Assert().True(nearQuestFound, "Near quest should be found within radius")
	s.Assert().False(farQuestFound, "Far quest should NOT be found outside radius")
}

func (s *Suite) TestSearchQuestsByRadiusHTTPEmpty() {
	ctx := context.Background()

	// Pre-condition - use location far from any existing quests
	remoteLocation := kernel.GeoCoordinate{
		Lat: -89.0, // Near South Pole
		Lon: 0.0,
	}

	// Act - search for quests in remote location via HTTP API
	radiusKm := 10.0
	searchReq := casesteps.SearchQuestsByRadiusHTTPRequest(
		remoteLocation.Latitude(), remoteLocation.Longitude(), radiusKm)
	searchResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, searchReq)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, searchResp.StatusCode)

	// Parse response
	var foundQuests []servers.Quest
	err = json.Unmarshal([]byte(searchResp.Body), &foundQuests)
	s.Require().NoError(err)

	s.Assert().Len(foundQuests, 0, "Should return empty list for remote location")
}

func (s *Suite) TestSearchQuestsByRadiusHTTPWithInvalidParams() {
	ctx := context.Background()

	// Test cases with invalid parameters
	testCases := []struct {
		name     string
		lat      float64
		lon      float64
		radiusKm float64
	}{
		{"Invalid latitude too high", 91.0, 0.0, 10.0},
		{"Invalid latitude too low", -91.0, 0.0, 10.0},
		{"Invalid longitude too high", 0.0, 181.0, 10.0},
		{"Invalid longitude too low", 0.0, -181.0, 10.0},
		{"Invalid radius zero", 50.0, 10.0, 0.0},
		{"Invalid radius negative", 50.0, 10.0, -5.0},
		{"Invalid radius too large", 50.0, 10.0, 25000.0},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act - send invalid parameters via HTTP API
			searchReq := casesteps.SearchQuestsByRadiusHTTPRequest(tc.lat, tc.lon, tc.radiusKm)
			searchResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, searchReq)

			// Assert - should return 400 validation error
			s.Require().NoError(err)
			s.Require().Equal(http.StatusBadRequest, searchResp.StatusCode, "Should return 400 for invalid parameters")

			// Verify error message contains validation details
			s.Assert().Contains(searchResp.Body, "validation failed", "Error message should contain validation failure details")
		})
	}
}
