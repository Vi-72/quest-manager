package http

import (
	"context"
	"github.com/google/uuid"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
)

// ChangeQuestStatus implements PATCH /api/v1/quests/{quest_id}/status from OpenAPI.
func (a *ApiHandler) ChangeQuestStatus(ctx context.Context, request servers.ChangeQuestStatusRequestObject) (servers.ChangeQuestStatusResponseObject, error) {
	if request.Body == nil {
		return servers.ChangeQuestStatus400Response{}, nil
	}

	questID, err := uuid.Parse(request.QuestId)
	if err != nil {
		return servers.ChangeQuestStatus404Response{}, nil
	}

	cmd := commands.ChangeQuestStatusCommand{
		ID:     questID,
		Status: quest.Status(request.Body.Status),
	}

	_, err = a.changeQuestStatusHandler.Handle(ctx, cmd)
	if err != nil {
		return servers.ChangeQuestStatus500Response{}, nil
	}

	// Get updated quest from repository
	updatedQuest, err := a.getQuestByIDHandler.Handle(ctx, queries.GetQuestByIDQuery{ID: questID})
	if err != nil {
		return servers.ChangeQuestStatus500Response{}, nil
	}

	apiQuest := QuestToAPI(updatedQuest.Quest)
	return servers.ChangeQuestStatus200JSONResponse(apiQuest), nil
}
