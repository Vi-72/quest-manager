package http

import (
	"context"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
)

// ListQuests implements GET /api/v1/quests from OpenAPI.
func (a *ApiHandler) ListQuests(ctx context.Context, request servers.ListQuestsRequestObject) (servers.ListQuestsResponseObject, error) {
	query := queries.ListQuestsQuery{}
	if request.Params.Status != nil {
		// Просто передаем статус как есть - домен/репозиторий сам разберется с валидностью
		statusStr := string(*request.Params.Status)
		status := quest.Status(statusStr)
		query.Status = &status
	}

	result, err := a.listQuestsHandler.Handle(ctx, query)
	if err != nil {
		return servers.ListQuests500Response{}, nil
	}

	var apiQuests []servers.Quest
	for _, q := range result.Quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.ListQuests200JSONResponse(apiQuests), nil
}
