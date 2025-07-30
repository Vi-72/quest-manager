package commands

import (
	"context"
	"fmt"

	"quest-manager/internal/core/ports"
)

// AssignQuestCommandHandler defines the interface for assigning a quest.
type AssignQuestCommandHandler interface {
	Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error)
}

var _ AssignQuestCommandHandler = &assignQuestHandler{}

type assignQuestHandler struct {
	repo ports.QuestRepository
}

// NewAssignQuestCommandHandler creates a new AssignQuestCommandHandler.
func NewAssignQuestCommandHandler(repo ports.QuestRepository) AssignQuestCommandHandler {
	return &assignQuestHandler{repo: repo}
}

// Handle assigns the quest to the specified user.
func (h *assignQuestHandler) Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error) {
	q, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return AssignQuestResult{}, fmt.Errorf("quest not found: %w", err)
	}

	if cmd.UserID == "" {
		return AssignQuestResult{}, fmt.Errorf("user_id is required for assignment")
	}

	q.AssignTo(cmd.UserID)

	if err := h.repo.Save(ctx, q); err != nil {
		return AssignQuestResult{}, fmt.Errorf("failed to save quest: %w", err)
	}

	return AssignQuestResult{
		ID:       q.ID,
		Assignee: cmd.UserID,
		Status:   string(q.Status),
	}, nil
}
