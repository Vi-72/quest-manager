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
	Handle(ctx context.Context, cmd CreateQuestCommand) (quest.Quest, error)
}

var _ CreateQuestCommandHandler = &createQuestHandler{}

type createQuestHandler struct {
	unitOfWork ports.UnitOfWork
}

// NewCreateQuestCommandHandler creates a new instance of CreateQuestCommandHandler.
func NewCreateQuestCommandHandler(unitOfWork ports.UnitOfWork) CreateQuestCommandHandler {
	return &createQuestHandler{unitOfWork: unitOfWork}
}

func (h *createQuestHandler) Handle(ctx context.Context, cmd CreateQuestCommand) (quest.Quest, error) {
	var targetLocationID *uuid.UUID
	var executionLocationID *uuid.UUID

	// Начинаем транзакцию
	h.unitOfWork.Begin(ctx)

	// Получаем адрес для target location
	var targetAddress string
	if cmd.TargetLocation.GetAddress() != nil {
		targetAddress = *cmd.TargetLocation.GetAddress()
	}

	// Всегда создаем локацию для target с пустым именем
	targetLoc, err := location.NewLocation(
		"", // пустое имя
		cmd.TargetLocation,
		targetAddress, // адрес из координат
		"",            // пустое описание
	)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, err
	}

	if err := h.unitOfWork.LocationRepository().Save(ctx, targetLoc); err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save target location", err)
	}

	locID := targetLoc.ID()
	targetLocationID = &locID

	// Получаем адрес для execution location
	var executionAddress string
	if cmd.ExecutionLocation.GetAddress() != nil {
		executionAddress = *cmd.ExecutionLocation.GetAddress()
	}

	// Всегда создаем локацию для execution с пустым именем
	executionLoc, err := location.NewLocation(
		"", // пустое имя
		cmd.ExecutionLocation,
		executionAddress, // адрес из координат
		"",               // пустое описание
	)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, err
	}

	if err := h.unitOfWork.LocationRepository().Save(ctx, executionLoc); err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save execution location", err)
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
		return quest.Quest{}, err
	}

	// Связываем квест с созданными локациями
	q.TargetLocationID = targetLocationID
	q.ExecutionLocationID = executionLocationID

	// Сохраняем квест
	if err := h.unitOfWork.QuestRepository().Save(ctx, q); err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Коммитим транзакцию
	if err := h.unitOfWork.Commit(ctx); err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to commit transaction", err)
	}

	return q, nil
}
