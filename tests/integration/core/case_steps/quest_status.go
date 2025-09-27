package casesteps

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// ChangeQuestStatusStep изменяет статус квеста
func ChangeQuestStatusStep(
	ctx context.Context,
	handler commands.ChangeQuestStatusCommandHandler,
	questRepo ports.QuestRepository,
	questID uuid.UUID,
	newStatus quest.Status,
) (quest.Quest, error) {
	cmd := commands.ChangeQuestStatusCommand{
		QuestID: questID,
		Status:  newStatus,
	}

	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		return quest.Quest{}, err
	}

	// Для тестов получаем полный квест
	return questRepo.GetByID(ctx, result.ID)
}

// AssignQuestStep выполняет назначение квеста пользователю (также изменяет статус)
func AssignQuestStep(
	ctx context.Context,
	handler commands.AssignQuestCommandHandler,
	questID uuid.UUID,
	userID uuid.UUID,
) (commands.AssignQuestResult, error) {
	cmd := commands.AssignQuestCommand{
		ID:     questID,
		UserID: userID,
	}

	return handler.Handle(ctx, cmd)
}
