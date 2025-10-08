package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/errors"
	"quest-manager/internal/adapters/in/http/middleware"
)

// ListAssignedQuests implements GET /api/v1/quests/assigned from OpenAPI.
func (a *ApiHandler) ListAssignedQuests(ctx context.Context, request v1.ListAssignedQuestsRequestObject) (v1.ListAssignedQuestsResponseObject, error) {
	// Get authenticated user ID from context (set by auth middleware)
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.NewBadRequest("user ID not found in context")
	}

	// Pass user ID from token to handler
	quests, err := a.listAssignedQuestsHandler.Handle(ctx, userID)
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
