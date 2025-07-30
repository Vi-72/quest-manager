package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/generated/servers"
)

// ListAssignedQuests implements GET /api/v1/quests/assigned from OpenAPI.
func (a *ApiHandler) ListAssignedQuests(ctx context.Context, request servers.ListAssignedQuestsRequestObject) (servers.ListAssignedQuestsResponseObject, error) {
	// Валидация UUID для user_id
	_, validationErr := validations.ValidateUUID(request.Params.UserId, "user_id")
	if validationErr != nil {
		// Возвращаем ошибку валидации, middleware автоматически обработает её и вернет 400 ответ
		return nil, validationErr
	}

	query := queries.ListAssignedQuestsQuery{UserID: request.Params.UserId}
	result, err := a.listAssignedQuestsHandler.Handle(ctx, query)
	if err != nil {
		return servers.ListAssignedQuests500Response{}, nil
	}

	var apiQuests []servers.Quest
	for _, q := range result.Quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.ListAssignedQuests200JSONResponse(apiQuests), nil
}
