package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/errors"
	"quest-manager/internal/core/domain/model/kernel"
)

// SearchQuestsByRadius implements GET /api/v1/quests/search-radius from OpenAPI.
func (a *ApiHandler) SearchQuestsByRadius(ctx context.Context, request v1.SearchQuestsByRadiusRequestObject) (v1.SearchQuestsByRadiusResponseObject, error) {
	center, err := kernel.NewGeoCoordinate(float64(request.Params.Lat), float64(request.Params.Lon))
	if err != nil {
		return nil, errors.NewBadRequest("Request validation failed: coordinates invalid (" + err.Error() + ")")
	}

	quests, err := a.searchQuestsByRadius.Handle(ctx, center, float64(request.Params.RadiusKm))
	if err != nil {
		// Pass error to middleware for proper handling
		return nil, err
	}

	var apiQuests []v1.Quest
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return v1.SearchQuestsByRadius200JSONResponse(apiQuests), nil
}
