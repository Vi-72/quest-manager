package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/problems"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
)

// SearchQuestsByRadius implements GET /api/v1/quests/search-radius from OpenAPI.
func (a *ApiHandler) SearchQuestsByRadius(ctx context.Context, request servers.SearchQuestsByRadiusRequestObject) (servers.SearchQuestsByRadiusResponseObject, error) {
	// Validate latitude range
	lat := float64(request.Params.Lat)
	if lat < -90 || lat > 90 {
		return nil, problems.NewBadRequest("Latitude must be between -90 and 90 degrees")
	}

	// Validate longitude range
	lon := float64(request.Params.Lon)
	if lon < -180 || lon > 180 {
		return nil, problems.NewBadRequest("Longitude must be between -180 and 180 degrees")
	}

	// Validate radius
	radiusKm := float64(request.Params.RadiusKm)
	if radiusKm <= 0 {
		return nil, problems.NewBadRequest("Radius must be greater than 0 kilometers")
	}
	if radiusKm > 20000 { // Earth's circumference is ~40000km, so max half of that
		return nil, problems.NewBadRequest("Radius must be less than 20000 kilometers")
	}

	center, err := kernel.NewGeoCoordinate(lat, lon)
	if err != nil {
		return nil, problems.NewBadRequest("Invalid coordinates: " + err.Error())
	}

	query := queries.SearchQuestsByRadiusQuery{
		Center:   center,
		RadiusKm: radiusKm,
	}

	result, err := a.searchQuestsByRadius.Handle(ctx, query)
	if err != nil {
		return servers.SearchQuestsByRadius500Response{}, nil
	}

	var apiQuests []servers.Quest
	for _, q := range result.Quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.SearchQuestsByRadius200JSONResponse(apiQuests), nil
}
