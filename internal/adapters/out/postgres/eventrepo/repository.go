package eventrepo

import (
	"context"
	"time"

	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"
)

var _ ports.EventPublisher = &Repository{}

type Repository struct {
	tracker ports.Tracker
}

func NewRepository(tracker ports.Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}
	return &Repository{tracker: tracker}, nil
}

// Publish сохраняет доменные события в базу данных
func (r *Repository) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	for _, event := range events {
		dto, err := r.domainEventToDTO(event)
		if err != nil {
			if !isInTransaction {
				_ = r.tracker.Rollback()
			}
			return errs.WrapInfrastructureError("failed to convert event to DTO", err)
		}

		err = tx.WithContext(ctx).Create(&dto).Error
		if err != nil {
			if !isInTransaction {
				_ = r.tracker.Rollback()
			}
			return errs.WrapInfrastructureError("failed to save event", err)
		}
	}

	if !isInTransaction {
		if err := r.tracker.Commit(ctx); err != nil {
			return errs.WrapInfrastructureError("failed to commit event transaction", err)
		}
	}

	return nil
}

// domainEventToDTO конвертирует доменное событие в DTO
func (r *Repository) domainEventToDTO(event ddd.DomainEvent) (EventDTO, error) {
	dto := EventDTO{
		ID:        event.GetID().String(),
		EventType: event.GetName(),
		CreatedAt: time.Now(),
	}

	// Определяем AggregateID и данные в зависимости от типа события
	switch e := event.(type) {
	// === QUEST EVENTS ===
	case quest.QuestCreated:
		dto.AggregateID = e.QuestID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	case quest.QuestStatusChanged:
		dto.AggregateID = e.QuestID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	case quest.QuestAssigned:
		dto.AggregateID = e.QuestID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	// === LOCATION EVENTS ===
	case location.LocationCreated:
		dto.AggregateID = e.Coordinate.LocationID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	case location.LocationUpdated:
		dto.AggregateID = e.Coordinate.LocationID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	default:
		// Fallback для неизвестных событий
		dto.AggregateID = event.GetID().String()
		data, err := MarshalEventData(event)
		if err != nil {
			return EventDTO{}, errs.NewDomainValidationError("eventSerialization", "failed to serialize unknown event type")
		}
		dto.Data = data
	}

	return dto, nil
}
