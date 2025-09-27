package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
)

// ListAssignedQuests implements GET /api/v1/quests/assigned from OpenAPI.
func (a *ApiHandler) ListAssignedQuests(ctx context.Context, request v1.ListAssignedQuestsRequestObject) (v1.ListAssignedQuestsResponseObject, error) {
	// UserId is already UUID type from OpenAPI, pass directly to handler
	quests, err := a.listAssignedQuestsHandler.Handle(ctx, request.Params.UserId)
	if err != nil {
		// Pass error to middleware for proper handling
		return nil, err
	}

	apiQuests := make([]v1.Quest, 0)
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return v1.ListAssignedQuests200JSONResponse(apiQuests), nil
}
