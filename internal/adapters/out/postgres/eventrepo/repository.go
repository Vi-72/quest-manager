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
	tracker            ports.Tracker
	goroutineSemaphore chan struct{} // Семафор для ограничения горутин
}

func NewRepository(tracker ports.Tracker, goroutineLimit int) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}
	if goroutineLimit <= 0 {
		goroutineLimit = 5 // значение по умолчанию
	}

	return &Repository{
		tracker:            tracker,
		goroutineSemaphore: make(chan struct{}, goroutineLimit),
	}, nil
}

// PublishAsync асинхронно публикует события с ограничением горутин
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

		// Публикуем события
		if err := r.Publish(ctx, events...); err != nil {
			// TODO: добавить логгер для записи ошибок
			_ = err
		}
	}()
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
		dto.AggregateID = e.AggregateID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	case quest.QuestStatusChanged:
		dto.AggregateID = e.AggregateID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	case quest.QuestAssigned:
		dto.AggregateID = e.AggregateID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	// === LOCATION EVENTS ===
	case location.LocationCreated:
		dto.AggregateID = e.AggregateID.String()
		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	case location.LocationUpdated:
		dto.AggregateID = e.AggregateID.String()
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
