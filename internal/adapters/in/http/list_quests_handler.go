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
		// Simply pass status as is - domain/repository will handle validity itself
		statusStr := string(*request.Params.Status)
		questStatus := quest.Status(statusStr)
		status = &questStatus
	}

	// Get quest list directly with optional filter
	quests, err := a.listQuestsHandler.Handle(ctx, status)
	if err != nil {
		// Pass error to middleware for proper handling (e.g., 400 for invalid status)
		return nil, err
	}

	var apiQuests []servers.Quest
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.ListQuests200JSONResponse(apiQuests), nil
}
