package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/problems"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
)

// ListAssignedQuests implements GET /api/v1/quests/assigned from OpenAPI.
func (a *ApiHandler) ListAssignedQuests(ctx context.Context, request servers.ListAssignedQuestsRequestObject) (servers.ListAssignedQuestsResponseObject, error) {
	if request.Params.UserId == "" {
		return nil, problems.NewBadRequest("UserId query parameter is required")
	}

	// Validate UserId is a valid UUID format
	_, err := uuid.Parse(request.Params.UserId)
	if err != nil {
		return nil, problems.NewBadRequest("UserId must be a valid UUID format (e.g. d5cde057-d462-419b-9428-42eebe22a85e)")
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
