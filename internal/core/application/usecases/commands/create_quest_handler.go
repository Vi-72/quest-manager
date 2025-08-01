package commands

import (
	"context"
	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"

	"github.com/google/uuid"
)

// CreateQuestCommandHandler defines a handler for creating a new quest.
type CreateQuestCommandHandler interface {
	Handle(ctx context.Context, cmd CreateQuestCommand) (quest.Quest, error)
}

var _ CreateQuestCommandHandler = &createQuestHandler{}

type createQuestHandler struct {
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
}

// NewCreateQuestCommandHandler creates a new instance of CreateQuestCommandHandler.
func NewCreateQuestCommandHandler(unitOfWork ports.UnitOfWork, eventPublisher ports.EventPublisher) CreateQuestCommandHandler {
	return &createQuestHandler{
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}
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

	// Создаем или находим target location
	targetLoc, err := location.NewLocation(
		cmd.TargetLocation,
		targetAddress,
	)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, errs.WrapInfrastructureError("failed to create target location", err)
	}

	// Сохраняем target location
	err = h.unitOfWork.LocationRepository().Save(ctx, targetLoc)
	if err != nil {
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

	// Создаем или находим execution location (может быть такой же как target)
	var executionLoc *location.Location
	if cmd.TargetLocation.Equals(cmd.ExecutionLocation) {
		executionLoc = targetLoc
		executionLocationID = targetLocationID
	} else {
		executionLoc, err = location.NewLocation(
			cmd.ExecutionLocation,
			executionAddress,
		)
		if err != nil {
			_ = h.unitOfWork.Rollback()
			return quest.Quest{}, errs.WrapInfrastructureError("failed to create execution location", err)
		}

		// Сохраняем execution location
		err = h.unitOfWork.LocationRepository().Save(ctx, executionLoc)
		if err != nil {
			_ = h.unitOfWork.Rollback()
			return quest.Quest{}, errs.WrapInfrastructureError("failed to save execution location", err)
		}
		locID = executionLoc.ID()
		executionLocationID = &locID
	}

	// Создаем квест
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
		return quest.Quest{}, errs.NewDomainValidationErrorWithCause("quest", "invalid quest data", err)
	}

	// Связываем квест с созданными локациями
	q.TargetLocationID = targetLocationID
	q.ExecutionLocationID = executionLocationID

	// Сохраняем квест
	err = h.unitOfWork.QuestRepository().Save(ctx, q)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Коммитим транзакцию
	err = h.unitOfWork.Commit(ctx)
	if err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to commit quest creation transaction", err)
	}

	// Публикуем все доменные события асинхронно после успешного коммита
	if h.eventPublisher != nil {
		var allEvents []ddd.DomainEvent

		// Добавляем события квеста
		allEvents = append(allEvents, q.GetDomainEvents()...)

		// Добавляем события target location
		allEvents = append(allEvents, targetLoc.GetDomainEvents()...)

		// Добавляем события execution location (если это не та же локация)
		if executionLoc != targetLoc {
			allEvents = append(allEvents, executionLoc.GetDomainEvents()...)
		}

		// Отправляем события в горутине
		go func() {
			if err := h.eventPublisher.Publish(context.Background(), allEvents...); err != nil {
				// Логируем ошибку, но не возвращаем её пользователю
				// TODO: добавить логгер для записи ошибок публикации событий
				_ = err
			}
		}()
	}

	// Очищаем события после постановки в очередь на публикацию
	q.ClearDomainEvents()
	targetLoc.ClearDomainEvents()
	if executionLoc != targetLoc {
		executionLoc.ClearDomainEvents()
	}

	return q, nil
}
