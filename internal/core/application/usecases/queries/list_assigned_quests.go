package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// ListAssignedQuestsQuery represents the input parameters to list quests assigned to a user.
type ListAssignedQuestsQuery struct {
	UserID string
}

// ListAssignedQuestsResult represents the result of listing quests assigned to a user.
type ListAssignedQuestsResult struct {
	Quests []quest.Quest
}

// ListAssignedQuestsQueryHandler defines the interface for handling ListAssignedQuestsQuery.
type ListAssignedQuestsQueryHandler interface {
	Handle(ctx context.Context, query ListAssignedQuestsQuery) (ListAssignedQuestsResult, error)
}

type listAssignedQuestsHandler struct {
	repo ports.QuestRepository
}

// NewListAssignedQuestsQueryHandler creates a new instance of ListAssignedQuestsQueryHandler.
func NewListAssignedQuestsQueryHandler(repo ports.QuestRepository) ListAssignedQuestsQueryHandler {
	return &listAssignedQuestsHandler{repo: repo}
}

// Handle retrieves all quests assigned to the given user.
func (h *listAssignedQuestsHandler) Handle(ctx context.Context, query ListAssignedQuestsQuery) (ListAssignedQuestsResult, error) {
	quests, err := h.repo.FindByAssignee(ctx, query.UserID)
	if err != nil {
		return ListAssignedQuestsResult{}, err
	}
	return ListAssignedQuestsResult{Quests: quests}, nil
}
