package eventrepo

import (
	"context"
	"log"
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

// PublishAsync asynchronously publishes events
func (r *Repository) PublishAsync(ctx context.Context, events ...ddd.DomainEvent) {
	if len(events) == 0 {
		return
	}

	// Run in goroutine for async behavior
	go func() {
		if err := r.Publish(context.Background(), events...); err != nil {
			log.Printf("ERROR: Failed to publish events: %v", err)
		}
	}()
}

// Publish сохраняет доменные события в базу данных
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

	// Save all events in a single transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range dtos {
			if err := tx.Create(&dtos[i]).Error; err != nil {
				return errs.WrapInfrastructureError("failed to save event", err)
			}
		}
		return nil
	})
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
