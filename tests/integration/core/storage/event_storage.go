package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quest-manager/internal/adapters/out/postgres/eventrepo"
)

// EventStorage предоставляет методы для работы с событиями в тестах
type EventStorage struct {
	db *gorm.DB
}

// NewEventStorage создает новый EventStorage
func NewEventStorage(db *gorm.DB) *EventStorage {
	return &EventStorage{db: db}
}

// GetEventByID получает событие из базы данных по ID
func (s *EventStorage) GetEventByID(ctx context.Context, id string) (*eventrepo.EventDTO, error) {
	var eventDTO eventrepo.EventDTO
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&eventDTO).Error
	if err != nil {
		return nil, err
	}
	return &eventDTO, nil
}

// GetAllEvents получает все события из базы данных
func (s *EventStorage) GetAllEvents(ctx context.Context) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	err := s.db.WithContext(ctx).Order("created_at ASC").Find(&events).Error
	return events, err
}

// GetEventsByType получает события по типу
func (s *EventStorage) GetEventsByType(ctx context.Context, eventType string) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	err := s.db.WithContext(ctx).Where("event_type = ?", eventType).Order("created_at ASC").Find(&events).Error
	return events, err
}

// GetEventsByAggregateID получает события по ID агрегата
func (s *EventStorage) GetEventsByAggregateID(ctx context.Context, aggregateID uuid.UUID) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	err := s.db.WithContext(ctx).Where("aggregate_id = ?", aggregateID.String()).Order("created_at ASC").Find(&events).Error
	return events, err
}

// GetEventsByTimeRange получает события в диапазоне времени
func (s *EventStorage) GetEventsByTimeRange(ctx context.Context, from, to time.Time) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	err := s.db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", from, to).Order("created_at ASC").Find(&events).Error
	return events, err
}

// GetQuestEvents получает все события связанные с квестами
func (s *EventStorage) GetQuestEvents(ctx context.Context) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	err := s.db.WithContext(ctx).Where("event_type LIKE 'quest.%'").Order("created_at ASC").Find(&events).Error
	return events, err
}

// GetLocationEvents получает все события связанные с локациями
func (s *EventStorage) GetLocationEvents(ctx context.Context) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	err := s.db.WithContext(ctx).Where("event_type LIKE 'location.%'").Order("created_at ASC").Find(&events).Error
	return events, err
}

// CountEvents подсчитывает количество событий
func (s *EventStorage) CountEvents(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&eventrepo.EventDTO{}).Count(&count).Error
	return count, err
}

// CountEventsByType подсчитывает количество событий по типу
func (s *EventStorage) CountEventsByType(ctx context.Context, eventType string) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&eventrepo.EventDTO{}).Where("event_type = ?", eventType).Count(&count).Error
	return count, err
}

// CountEventsByAggregateID подсчитывает количество событий для агрегата
func (s *EventStorage) CountEventsByAggregateID(ctx context.Context, aggregateID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&eventrepo.EventDTO{}).Where("aggregate_id = ?", aggregateID.String()).Count(&count).Error
	return count, err
}

// WaitForEvents ждет появления определенного количества событий (для асинхронных операций)
func (s *EventStorage) WaitForEvents(ctx context.Context, expectedCount int64, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		count, err := s.CountEvents(ctx)
		if err != nil {
			return err
		}

		if count >= expectedCount {
			return nil
		}

		time.Sleep(10 * time.Millisecond)
	}

	return ErrTimeout
}

// WaitForEventsOfType ждет появления событий определенного типа
func (s *EventStorage) WaitForEventsOfType(ctx context.Context, eventType string, expectedCount int64, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		count, err := s.CountEventsByType(ctx, eventType)
		if err != nil {
			return err
		}

		if count >= expectedCount {
			return nil
		}

		time.Sleep(10 * time.Millisecond)
	}

	return ErrTimeout
}

// DeleteEventByID удаляет событие по ID
func (s *EventStorage) DeleteEventByID(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(&eventrepo.EventDTO{}, "id = ?", id).Error
}

// DeleteAllEvents удаляет все события
func (s *EventStorage) DeleteAllEvents(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec("DELETE FROM events").Error
}

// ErrTimeout ошибка таймаута ожидания
var ErrTimeout = gorm.ErrRecordNotFound
