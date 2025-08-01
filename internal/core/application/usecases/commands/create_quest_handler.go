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

	// Begin transaction
	h.unitOfWork.Begin(ctx)

	// Create or find target location
	targetLoc, err := location.NewLocation(
		cmd.TargetLocation,
		cmd.TargetAddress,
	)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, errs.WrapInfrastructureError("failed to create target location", err)
	}

	// Save target location
	err = h.unitOfWork.LocationRepository().Save(ctx, targetLoc)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save target location", err)
	}
	locID := targetLoc.ID()
	targetLocationID = &locID

	// Create or find execution location (can be the same as target)
	var executionLoc *location.Location
	if cmd.TargetLocation.Equals(cmd.ExecutionLocation) {
		executionLoc = targetLoc
		executionLocationID = targetLocationID
	} else {
		executionLoc, err = location.NewLocation(
			cmd.ExecutionLocation,
			cmd.ExecutionAddress,
		)
		if err != nil {
			_ = h.unitOfWork.Rollback()
			return quest.Quest{}, errs.WrapInfrastructureError("failed to create execution location", err)
		}

		// Save execution location
		err = h.unitOfWork.LocationRepository().Save(ctx, executionLoc)
		if err != nil {
			_ = h.unitOfWork.Rollback()
			return quest.Quest{}, errs.WrapInfrastructureError("failed to save execution location", err)
		}
		locID = executionLoc.ID()
		executionLocationID = &locID
	}

	// Create quest
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

	// Link quest with created locations
	q.TargetLocationID = targetLocationID
	q.ExecutionLocationID = executionLocationID

	// Save quest
	err = h.unitOfWork.QuestRepository().Save(ctx, q)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Commit transaction
	err = h.unitOfWork.Commit(ctx)
	if err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to commit quest creation transaction", err)
	}

	// Publish all domain events asynchronously after successful commit
	if h.eventPublisher != nil {
		var allEvents []ddd.DomainEvent

		// Add quest events
		allEvents = append(allEvents, q.GetDomainEvents()...)

		// Add target location events
		allEvents = append(allEvents, targetLoc.GetDomainEvents()...)

		// Add execution location events (if it's not the same location)
		if executionLoc != targetLoc {
			allEvents = append(allEvents, executionLoc.GetDomainEvents()...)
		}

		// Send events asynchronously with goroutine limiting
		h.eventPublisher.PublishAsync(context.Background(), allEvents...)
	}

	// Clear events after queuing for publication
	q.ClearDomainEvents()
	targetLoc.ClearDomainEvents()
	if executionLoc != targetLoc {
		executionLoc.ClearDomainEvents()
	}

	return q, nil
}
