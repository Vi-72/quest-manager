package storage

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quest-manager/internal/adapters/out/postgres/locationrepo"
)

// LocationStorage предоставляет методы для работы с локациями в тестах
type LocationStorage struct {
	db *gorm.DB
}

// NewLocationStorage создает новый LocationStorage
func NewLocationStorage(db *gorm.DB) *LocationStorage {
	return &LocationStorage{db: db}
}

// GetLocationByID получает локацию из базы данных по ID
func (s *LocationStorage) GetLocationByID(ctx context.Context, id uuid.UUID) (*locationrepo.LocationDTO, error) {
	var locationDTO locationrepo.LocationDTO
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&locationDTO).Error
	if err != nil {
		return nil, err
	}
	return &locationDTO, nil
}

// GetAllLocations получает все локации из базы данных
func (s *LocationStorage) GetAllLocations(ctx context.Context) ([]locationrepo.LocationDTO, error) {
	var locations []locationrepo.LocationDTO
	err := s.db.WithContext(ctx).Find(&locations).Error
	return locations, err
}

// GetLocationsByAddress получает локации по адресу (частичное совпадение)
func (s *LocationStorage) GetLocationsByAddress(ctx context.Context, address string) ([]locationrepo.LocationDTO, error) {
	var locations []locationrepo.LocationDTO
	err := s.db.WithContext(ctx).Where("address ILIKE ?", "%"+address+"%").Find(&locations).Error
	return locations, err
}

// GetLocationsInRadius получает локации в радиусе от точки
func (s *LocationStorage) GetLocationsInRadius(ctx context.Context, lat, lon, radiusKm float64) ([]locationrepo.LocationDTO, error) {
	var locations []locationrepo.LocationDTO
	// Простая проверка по квадрату (для тестов достаточно)
	latRange := radiusKm / 111.0         // Примерно 1 градус = 111 км
	lonRange := radiusKm / (111.0 * 0.7) // Учитываем широту

	err := s.db.WithContext(ctx).Where(
		"latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
		lat-latRange, lat+latRange, lon-lonRange, lon+lonRange,
	).Find(&locations).Error
	return locations, err
}

// CountLocations подсчитывает количество локаций
func (s *LocationStorage) CountLocations(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&locationrepo.LocationDTO{}).Count(&count).Error
	return count, err
}

// DeleteLocationByID удаляет локацию по ID
func (s *LocationStorage) DeleteLocationByID(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&locationrepo.LocationDTO{}, "id = ?", id).Error
}

// DeleteAllLocations удаляет все локации
func (s *LocationStorage) DeleteAllLocations(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec("DELETE FROM locations").Error
}
