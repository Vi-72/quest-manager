package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
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

	// Получаем список квестов напрямую
	quests, err := a.listAssignedQuestsHandler.Handle(ctx, request.Params.UserId)
	if err != nil {
		// Передаем ошибку в middleware для правильной обработки
		return nil, err
	}

	var apiQuests []servers.Quest
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.ListAssignedQuests200JSONResponse(apiQuests), nil
}
