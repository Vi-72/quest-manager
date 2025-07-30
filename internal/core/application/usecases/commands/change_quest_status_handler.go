package commands

import (
	"context"
	"fmt"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// ChangeQuestStatusCommandHandler defines the interface for handling ChangeQuestStatusCommand.
type ChangeQuestStatusCommandHandler interface {
	Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error)
}

var _ ChangeQuestStatusCommandHandler = &changeQuestStatusHandler{}

// changeQuestStatusHandler implements ChangeQuestStatusCommandHandler.
type changeQuestStatusHandler struct {
	repo ports.QuestRepository
}

// NewChangeQuestStatusCommandHandler creates a new instance of ChangeQuestStatusCommandHandler.
func NewChangeQuestStatusCommandHandler(repo ports.QuestRepository) ChangeQuestStatusCommandHandler {
	return &changeQuestStatusHandler{repo: repo}
}

// Handle updates the quest status using domain business rules.
func (h *changeQuestStatusHandler) Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error) {
	q, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return ChangeQuestStatusResult{}, fmt.Errorf("quest not found: %w", err)
	}

	// Используем доменную логику вместо switch/case
	if err := q.ChangeStatus(quest.Status(cmd.Status)); err != nil {
		return ChangeQuestStatusResult{}, fmt.Errorf("failed to change quest status: %w", err)
	}

	if err := h.repo.Save(ctx, q); err != nil {
		return ChangeQuestStatusResult{}, fmt.Errorf("failed to save quest: %w", err)
	}

	return ChangeQuestStatusResult{
		ID:     q.ID,
		Status: q.Status,
	}, nil
}
