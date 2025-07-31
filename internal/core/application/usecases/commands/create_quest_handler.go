package commands

import (
	"context"
	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"

	"github.com/google/uuid"
)

// CreateQuestCommandHandler defines a handler for creating a new quest.
type CreateQuestCommandHandler interface {
	Handle(ctx context.Context, cmd CreateQuestCommand) (CreateQuestResult, error)
}

var _ CreateQuestCommandHandler = &createQuestHandler{}

type createQuestHandler struct {
	unitOfWork ports.UnitOfWork
}

// NewCreateQuestCommandHandler creates a new instance of CreateQuestCommandHandler.
func NewCreateQuestCommandHandler(unitOfWork ports.UnitOfWork) CreateQuestCommandHandler {
	return &createQuestHandler{unitOfWork: unitOfWork}
}

func (h *createQuestHandler) Handle(ctx context.Context, cmd CreateQuestCommand) (CreateQuestResult, error) {
	var targetLocationID *uuid.UUID
	var executionLocationID *uuid.UUID

	// Начинаем транзакцию
	h.unitOfWork.Begin(ctx)

	// Всегда создаем локацию для target с пустым именем
	targetLoc, err := location.NewLocation(
		"", // пустое имя
		cmd.TargetLocation,
		"", // пустой адрес
		"", // пустое описание
	)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return CreateQuestResult{}, err
	}

	if err := h.unitOfWork.LocationRepository().Save(ctx, targetLoc); err != nil {
		_ = h.unitOfWork.Rollback()
		return CreateQuestResult{}, errs.WrapInfrastructureError("failed to save target location", err)
	}

	locID := targetLoc.ID()
	targetLocationID = &locID

	// Всегда создаем локацию для execution с пустым именем
	executionLoc, err := location.NewLocation(
		"", // пустое имя
		cmd.ExecutionLocation,
		"", // пустой адрес
		"", // пустое описание
	)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return CreateQuestResult{}, err
	}

	if err := h.unitOfWork.LocationRepository().Save(ctx, executionLoc); err != nil {
		_ = h.unitOfWork.Rollback()
		return CreateQuestResult{}, errs.WrapInfrastructureError("failed to save execution location", err)
	}

	locID = executionLoc.ID()
	executionLocationID = &locID

	// Создаем новый квест
	q, err := quest.NewQuest(
		cmd.Title,
		cmd.Description,
		cmd.Difficulty,
		cmd.Reward,
		cmd.DurationMinutes,
		cmd.TargetLocation,
		cmd.ExecutionLocation,
		cmd.Creator,
		cmd.Equipment,
		cmd.Skills,
	)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return CreateQuestResult{}, err
	}

	// Связываем квест с созданными локациями
	q.TargetLocationID = targetLocationID
	q.ExecutionLocationID = executionLocationID

	// Сохраняем квест
	if err := h.unitOfWork.QuestRepository().Save(ctx, q); err != nil {
		_ = h.unitOfWork.Rollback()
		return CreateQuestResult{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Коммитим транзакцию
	if err := h.unitOfWork.Commit(ctx); err != nil {
		return CreateQuestResult{}, errs.WrapInfrastructureError("failed to commit transaction", err)
	}

	return CreateQuestResult{
		ID:                  q.ID(),
		CreatedAt:           q.CreatedAt,
		Status:              q.Status,
		TargetLocationID:    targetLocationID,
		ExecutionLocationID: executionLocationID,
	}, nil
}
