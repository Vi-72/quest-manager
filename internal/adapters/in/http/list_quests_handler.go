package http

import (
	"context"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
)

// ListQuests implements GET /api/v1/quests from OpenAPI.
func (a *ApiHandler) ListQuests(ctx context.Context, request servers.ListQuestsRequestObject) (servers.ListQuestsResponseObject, error) {
	var status *quest.Status
	if request.Params.Status != nil {
		// Просто передаем статус как есть - домен/репозиторий сам разберется с валидностью
		statusStr := string(*request.Params.Status)
		questStatus := quest.Status(statusStr)
		status = &questStatus
	}

	// Получаем список квестов напрямую с опциональным фильтром
	quests, err := a.listQuestsHandler.Handle(ctx, status)
	if err != nil {
		// Передаем ошибку в middleware для правильной обработки (например, 400 для невалидного статуса)
		return nil, err
	}

	var apiQuests []servers.Quest
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.ListQuests200JSONResponse(apiQuests), nil
}
