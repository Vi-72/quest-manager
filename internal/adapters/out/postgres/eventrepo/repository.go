package eventrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"
)

var _ ports.EventPublisher = &Repository{}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Publish сохраняет доменные события в базу данных
// Использует существующую транзакцию из r.db, если она есть
// (например, когда вызывается из TransactionManager.RunInTransaction)
func (r *Repository) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	// Convert all events to DTOs
	var dtos []EventDTO
	for _, event := range events {
		dto, err := r.domainEventToDTO(event)
		if err != nil {
			return errs.WrapInfrastructureError("failed to convert event to DTO", err)
		}
		dtos = append(dtos, dto)
	}

	// Use r.db directly - it's already a transaction when called from TransactionManager
	// If it's not a transaction, GORM will execute operations without transaction
	// This ensures events are part of the same transaction as domain changes
	db := r.db.WithContext(ctx)
	for i := range dtos {
		if err := db.Create(&dtos[i]).Error; err != nil {
			return errs.WrapInfrastructureError("failed to save event", err)
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
	// Обрабатываем явно поддерживаемые типы
	case quest.QuestCreated,
		quest.QuestStatusChanged,
		quest.QuestAssigned,
		location.LocationCreated,
		location.LocationUpdated:

		// Приводим к общему интерфейсу
		agg, ok := e.(interface {
			GetAggregateID() uuid.UUID
		})
		if !ok {
			return EventDTO{}, errs.NewDomainValidationError("eventSerialization", "event missing AggregateID")
		}

		dto.AggregateID = agg.GetAggregateID().String()

		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	default:
		// Fallback для неизвестных событий
		// Попробуем получить AggregateID через интерфейс, если есть
		if agg, ok := e.(interface {
			GetAggregateID() uuid.UUID
		}); ok {
			dto.AggregateID = agg.GetAggregateID().String()
		} else {
			// Если GetAggregateID недоступен, используем ID события как fallback
			dto.AggregateID = event.GetID().String()
		}
		data, err := MarshalEventData(event)
		if err != nil {
			return EventDTO{}, errs.NewDomainValidationError("eventSerialization", "failed to serialize unknown event type")
		}
		dto.Data = data
	}

	return dto, nil
}
