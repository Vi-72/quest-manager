package http

import (
	"context"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
)

// SearchQuestsByRadius implements GET /api/v1/quests/search-radius from OpenAPI.
func (a *ApiHandler) SearchQuestsByRadius(ctx context.Context, request servers.SearchQuestsByRadiusRequestObject) (servers.SearchQuestsByRadiusResponseObject, error) {
	center, err := kernel.NewGeoCoordinate(float64(request.Params.Lat), float64(request.Params.Lon))
	if err != nil {
		return servers.SearchQuestsByRadius400Response{}, nil
	}

	query := queries.SearchQuestsByRadiusQuery{
		Center:   center,
		RadiusKm: float64(request.Params.RadiusKm),
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
