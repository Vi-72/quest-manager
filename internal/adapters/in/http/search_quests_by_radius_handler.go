package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/generated/servers"
)

// SearchQuestsByRadius implements GET /api/v1/quests/search-radius from OpenAPI.
func (a *ApiHandler) SearchQuestsByRadius(ctx context.Context, request servers.SearchQuestsByRadiusRequestObject) (servers.SearchQuestsByRadiusResponseObject, error) {
	// Валидация параметров поиска
	validatedData, validationErr := validations.ValidateSearchByRadiusParams(
		request.Params.Lat,
		request.Params.Lon,
		request.Params.RadiusKm,
	)
	if validationErr != nil {
		// Возвращаем ошибку валидации, middleware автоматически обработает её и вернет 400 ответ
		return nil, validationErr
	}

	query := queries.SearchQuestsByRadiusQuery{
		Center:   validatedData.Center,
		RadiusKm: validatedData.RadiusKm,
	}

	// Получаем список квестов напрямую
	quests, err := a.searchQuestsByRadius.Handle(ctx, query)
	if err != nil {
		return servers.SearchQuestsByRadius500Response{}, nil
	}

	var apiQuests []servers.Quest
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.SearchQuestsByRadius200JSONResponse(apiQuests), nil
}
