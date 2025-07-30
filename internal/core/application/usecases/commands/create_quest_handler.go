package commands

import (
	"context"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// CreateQuestCommandHandler defines a handler for creating a new quest.
type CreateQuestCommandHandler interface {
	Handle(ctx context.Context, cmd CreateQuestCommand) (CreateQuestResult, error)
}

var _ CreateQuestCommandHandler = &createQuestHandler{}

type createQuestHandler struct {
	repo ports.QuestRepository
}

// NewCreateQuestCommandHandler creates a new instance of CreateQuestCommandHandler.
func NewCreateQuestCommandHandler(repo ports.QuestRepository) CreateQuestCommandHandler {
	return &createQuestHandler{repo: repo}
}

func (h *createQuestHandler) Handle(ctx context.Context, cmd CreateQuestCommand) (CreateQuestResult, error) {
	q := quest.NewQuest(
		cmd.Title,
		cmd.Description,
		cmd.Difficulty,
		cmd.Reward,
		cmd.TargetLocation,
		cmd.ExecutionLocation,
		cmd.Creator,
		cmd.Equipment,
		cmd.Skills,
	)

	if err := h.repo.Save(ctx, q); err != nil {
		return CreateQuestResult{}, err
	}

	return CreateQuestResult{
		ID:        q.ID,
		CreatedAt: q.CreatedAt,
		Status:    q.Status,
	}, nil
}
