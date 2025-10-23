package eventrepo

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"
)

var _ ports.EventPublisher = &Repository{}

type Repository struct {
	uowFactory         ports.UnitOfWorkFactory
	goroutineSemaphore chan struct{}
}

func NewRepository(uowFactory ports.UnitOfWorkFactory, goroutineLimit int) (*Repository, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("uowFactory")
	}
	if goroutineLimit <= 0 {
		goroutineLimit = 5 // default value
	}

	return &Repository{
		uowFactory:         uowFactory,
		goroutineSemaphore: make(chan struct{}, goroutineLimit),
	}, nil
}

// PublishAsync asynchronously publishes events with goroutine limiting
func (r *Repository) PublishAsync(ctx context.Context, events ...ddd.DomainEvent) {
	if len(events) == 0 {
		return
	}

	// Запускаем в горутине с ограничением
	go func() {
		// Занимаем слот в семафоре
		r.goroutineSemaphore <- struct{}{}
		defer func() {
			// Освобождаем слот
			<-r.goroutineSemaphore
		}()

		// Create new UoW for this async operation
		uow, err := r.uowFactory.CreateUnitOfWork()
		if err != nil {
			log.Printf("ERROR: Failed to create UoW for event publishing: %v", err)
			return
		}

		if err := r.publishWithUnitOfWork(ctx, uow, events...); err != nil {
			log.Printf("ERROR: Failed to publish events: %v", err)
		}
	}()
}

// Publish сохраняет доменные события в базу данных
func (r *Repository) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	// Prefer using existing UnitOfWork from context when available
	if existingUow, ok := ports.UoWFromCtx(ctx); ok {
		return r.publishWithUnitOfWork(ctx, existingUow, events...)
	}

	// Fallback: create a new UnitOfWork for synchronous publishing
	uow, err := r.uowFactory.CreateUnitOfWork()
	if err != nil {
		return errs.WrapInfrastructureError("failed to create unit of work for event publishing", err)
	}
	return r.publishWithUnitOfWork(ctx, uow, events...)
}

func (r *Repository) publishWithUnitOfWork(ctx context.Context, uow ports.UnitOfWork, events ...ddd.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	// Detect whether we are already in a transaction and manage tx boundaries accordingly
	tracker := uow.(ports.Tracker)
	alreadyInTx := tracker.InTx()
	if !alreadyInTx {
		if err := uow.Begin(ctx); err != nil {
			return errs.WrapInfrastructureError("failed to begin event transaction", err)
		}
	}
	tx := tracker.Tx()

	for _, event := range events {
		dto, err := r.domainEventToDTO(event)
		if err != nil {
			if !alreadyInTx {
				_ = uow.Rollback()
			}
			return errs.WrapInfrastructureError("failed to convert event to DTO", err)
		}

		err = tx.WithContext(ctx).Create(&dto).Error
		if err != nil {
			if !alreadyInTx {
				_ = uow.Rollback()
			}
			return errs.WrapInfrastructureError("failed to save event", err)
		}
	}

	if !alreadyInTx {
		if err := uow.Commit(ctx); err != nil {
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
