package eventrepo

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	postgres "quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"
)

var _ ports.EventPublisher = &Repository{}

type Repository struct {
	trackerFactory     func() (ports.Tracker, error)
	tracker            ports.Tracker
	goroutineSemaphore chan struct{} // Semaphore for limiting goroutines
	mu                 sync.Mutex
}

func NewRepository(tracker ports.Tracker, goroutineLimit int) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}
	if goroutineLimit <= 0 {
		goroutineLimit = 5 // default value
	}

	db := tracker.Db()

	return &Repository{
		tracker: tracker,
		trackerFactory: func() (ports.Tracker, error) {
			uow, err := postgres.NewUnitOfWork(db)
			if err != nil {
				return nil, err
			}
			return uow.(ports.Tracker), nil
		},
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

		tracker, err := r.trackerFactory()
		if err != nil {
			log.Printf("ERROR: Failed to create tracker for event publishing: %v", err)
			return
		}

		if err := r.publishWithTracker(ctx, tracker, events...); err != nil {
			log.Printf("ERROR: Failed to publish events: %v", err)
		}
	}()
}

// Publish сохраняет доменные события в базу данных
func (r *Repository) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.publishWithTracker(ctx, r.tracker, events...)
}

func (r *Repository) publishWithTracker(ctx context.Context, tracker ports.Tracker, events ...ddd.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	isInTransaction := tracker.InTx()
	if !isInTransaction {
		if err := tracker.Begin(ctx); err != nil {
			return errs.WrapInfrastructureError("failed to begin event transaction", err)
		}
	}
	tx := tracker.Tx()

	for _, event := range events {
		dto, err := r.domainEventToDTO(event)
		if err != nil {
			if !isInTransaction {
				_ = tracker.Rollback()
			}
			return errs.WrapInfrastructureError("failed to convert event to DTO", err)
		}

		err = tx.WithContext(ctx).Create(&dto).Error
		if err != nil {
			if !isInTransaction {
				_ = tracker.Rollback()
			}
			return errs.WrapInfrastructureError("failed to save event", err)
		}
	}

	if !isInTransaction {
		if err := tracker.Commit(ctx); err != nil {
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
		dto.AggregateID = event.GetID().String()
		data, err := MarshalEventData(event)
		if err != nil {
			return EventDTO{}, errs.NewDomainValidationError("eventSerialization", "failed to serialize unknown event type")
		}
		dto.Data = data
	}

	return dto, nil
}
