package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/generated/servers"
)

// GetQuestById implements GET /api/v1/quests/{quest_id} from OpenAPI.
func (a *ApiHandler) GetQuestById(ctx context.Context, request servers.GetQuestByIdRequestObject) (servers.GetQuestByIdResponseObject, error) {
	// Валидация UUID
	questID, validationErr := validations.ValidateUUID(request.QuestId, "quest_id")
	if validationErr != nil {
		// Возвращаем ошибку валидации, middleware автоматически обработает её и вернет 400 ответ
		return nil, validationErr
	}

	// Получаем квест напрямую
	quest, err := a.getQuestByIDHandler.Handle(ctx, questID)
	if err != nil {
		return servers.GetQuestById404Response{}, nil
	}

	apiQuest := QuestToAPI(quest)
	return servers.GetQuestById200JSONResponse(apiQuest), nil
}
