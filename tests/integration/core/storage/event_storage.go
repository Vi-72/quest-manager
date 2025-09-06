package storage

import (
	"context"
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

// CountEvents подсчитывает количество событий
func (s *EventStorage) CountEvents(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&eventrepo.EventDTO{}).Count(&count).Error
	return count, err
}
