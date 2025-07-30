package commands

import (
	"context"
	"fmt"

	"quest-manager/internal/core/ports"
)

// AssignQuestCommandHandler defines the interface for handling AssignQuestCommand.
type AssignQuestCommandHandler interface {
	Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error)
}

var _ AssignQuestCommandHandler = &assignQuestHandler{}

// assignQuestHandler implements AssignQuestCommandHandler.
type assignQuestHandler struct {
	repo ports.QuestRepository
}

// NewAssignQuestCommandHandler creates a new instance of AssignQuestCommandHandler.
func NewAssignQuestCommandHandler(repo ports.QuestRepository) AssignQuestCommandHandler {
	return &assignQuestHandler{repo: repo}
}

// Handle assigns a quest to a user using domain business rules.
func (h *assignQuestHandler) Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error) {
	q, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return AssignQuestResult{}, fmt.Errorf("quest not found: %w", err)
	}

	// Используем доменную логику вместо прямого изменения полей
	if err := q.AssignTo(cmd.UserID); err != nil {
		return AssignQuestResult{}, fmt.Errorf("failed to assign quest: %w", err)
	}

	if err := h.repo.Save(ctx, q); err != nil {
		return AssignQuestResult{}, fmt.Errorf("failed to save quest: %w", err)
	}

	return AssignQuestResult{
		ID:       q.ID,
		Assignee: cmd.UserID,
		Status:   string(q.Status),
	}, nil
}
