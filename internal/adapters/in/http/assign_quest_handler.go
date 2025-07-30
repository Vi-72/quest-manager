package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/problems"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
)

// AssignQuest implements POST /api/v1/quests/{quest_id}/assign from OpenAPI.
func (a *ApiHandler) AssignQuest(ctx context.Context, request servers.AssignQuestRequestObject) (servers.AssignQuestResponseObject, error) {
	if request.Body == nil {
		return nil, problems.NewBadRequest("Request body is required")
	}

	if request.Body.UserId == "" {
		return nil, problems.NewBadRequest("UserId is required")
	}

	// Validate UserId is a valid UUID format
	_, err := uuid.Parse(request.Body.UserId)
	if err != nil {
		return nil, problems.NewBadRequest("UserId must be a valid UUID format (e.g. d5cde057-d462-419b-9428-42eebe22a85e)")
	}

	questID, err := uuid.Parse(request.QuestId)
	if err != nil {
		return nil, problems.NewBadRequest("Quest ID must be a valid UUID format")
	}

	cmd := commands.AssignQuestCommand{
		ID:     questID,
		UserID: request.Body.UserId,
	}

	_, err = a.assignQuestHandler.Handle(ctx, cmd)
	if err != nil {
		return servers.AssignQuest500Response{}, nil
	}

	// Get updated quest from repository
	updatedQuest, err := a.getQuestByIDHandler.Handle(ctx, queries.GetQuestByIDQuery{ID: questID})
	if err != nil {
		return nil, problems.NewNotFound("Quest not found")
	}

	apiQuest := QuestToAPI(updatedQuest.Quest)
	return servers.AssignQuest200JSONResponse(apiQuest), nil
}
