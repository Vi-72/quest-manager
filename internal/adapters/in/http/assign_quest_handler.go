package http

import (
	"context"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/generated/servers"

	"github.com/google/uuid"
)

// AssignQuest implements POST /api/v1/quests/{quest_id}/assign from OpenAPI.
func (a *ApiHandler) AssignQuest(ctx context.Context, request servers.AssignQuestRequestObject) (servers.AssignQuestResponseObject, error) {
	if request.Body == nil || request.Body.UserId == "" {
		return servers.AssignQuest400Response{}, nil
	}

	questID, err := uuid.Parse(request.QuestId)
	if err != nil {
		return servers.AssignQuest404Response{}, nil
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
		return servers.AssignQuest500Response{}, nil
	}

	apiQuest := QuestToAPI(updatedQuest.Quest)
	return servers.AssignQuest200JSONResponse(apiQuest), nil
}
