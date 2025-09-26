package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/core/domain/model/quest"
)

// ListQuests implements GET /api/v1/quests from OpenAPI.
func (a *ApiHandler) ListQuests(ctx context.Context, request v1.ListQuestsRequestObject) (v1.ListQuestsResponseObject, error) {
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

	var apiQuests []v1.Quest
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return v1.ListQuests200JSONResponse(apiQuests), nil
}
