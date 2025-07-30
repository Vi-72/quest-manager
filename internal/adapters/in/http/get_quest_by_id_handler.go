package http

import (
	"context"
	"github.com/google/uuid"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/generated/servers"
)

// GetQuestById implements GET /api/v1/quests/{quest_id} from OpenAPI.
func (a *ApiHandler) GetQuestById(ctx context.Context, request servers.GetQuestByIdRequestObject) (servers.GetQuestByIdResponseObject, error) {
	questID, err := uuid.Parse(request.QuestId)
	if err != nil {
		return servers.GetQuestById404Response{}, nil
	}

	query := queries.GetQuestByIDQuery{ID: questID}
	result, err := a.getQuestByIDHandler.Handle(ctx, query)
	if err != nil {
		return servers.GetQuestById404Response{}, nil
	}

	apiQuest := QuestToAPI(result.Quest)
	return servers.GetQuestById200JSONResponse(apiQuest), nil
}
