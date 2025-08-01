package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/generated/servers"
)

// SearchQuestsByRadius implements GET /api/v1/quests/search-radius from OpenAPI.
func (a *ApiHandler) SearchQuestsByRadius(ctx context.Context, request servers.SearchQuestsByRadiusRequestObject) (servers.SearchQuestsByRadiusResponseObject, error) {
	// Validate search parameters
	validatedData, validationErr := validations.ValidateSearchByRadiusParams(
		request.Params.Lat,
		request.Params.Lon,
		request.Params.RadiusKm,
	)
	if validationErr != nil {
		// Return validation error, middleware will automatically handle it and return 400 response
		return nil, validationErr
	}

	// Get quest list directly
	quests, err := a.searchQuestsByRadius.Handle(ctx, validatedData.Center, validatedData.RadiusKm)
	if err != nil {
		// Pass error to middleware for proper handling
		return nil, err
	}

	var apiQuests []servers.Quest
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.SearchQuestsByRadius200JSONResponse(apiQuests), nil
}
