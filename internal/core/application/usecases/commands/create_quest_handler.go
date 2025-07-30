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
	// Создаем новый квест (домен теперь валидирует difficulty)
	q, err := quest.NewQuest(
		cmd.Title,
		cmd.Description,
		cmd.Difficulty, // Теперь уже string, не нужна конвертация
		cmd.Reward,
		cmd.TargetLocation,
		cmd.ExecutionLocation,
		cmd.Creator,
		cmd.Equipment,
		cmd.Skills,
	)
	if err != nil {
		return CreateQuestResult{}, err
	}

	// Сохраняем квест
	if err := h.repo.Save(ctx, q); err != nil {
		return CreateQuestResult{}, err
	}

	return CreateQuestResult{
		ID:        q.ID,
		CreatedAt: q.CreatedAt,
		Status:    q.Status,
	}, nil
}
