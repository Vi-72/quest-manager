package storage

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quest-manager/internal/adapters/out/postgres/questrepo"
	"quest-manager/internal/core/domain/model/quest"
)

// QuestStorage предоставляет методы для работы с квестами в тестах
type QuestStorage struct {
	db *gorm.DB
}

// NewQuestStorage создает новый QuestStorage
func NewQuestStorage(db *gorm.DB) *QuestStorage {
	return &QuestStorage{db: db}
}

// GetQuestByID получает квест из базы данных по ID
func (s *QuestStorage) GetQuestByID(ctx context.Context, id uuid.UUID) (*questrepo.QuestDTO, error) {
	var questDTO questrepo.QuestDTO
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&questDTO).Error
	if err != nil {
		return nil, err
	}
	return &questDTO, nil
}

// GetAllQuests получает все квесты из базы данных
func (s *QuestStorage) GetAllQuests(ctx context.Context) ([]questrepo.QuestDTO, error) {
	var quests []questrepo.QuestDTO
	err := s.db.WithContext(ctx).Find(&quests).Error
	return quests, err
}

// GetQuestsByStatus получает квесты по статусу
func (s *QuestStorage) GetQuestsByStatus(ctx context.Context, status quest.Status) ([]questrepo.QuestDTO, error) {
	var quests []questrepo.QuestDTO
	err := s.db.WithContext(ctx).Where("status = ?", string(status)).Find(&quests).Error
	return quests, err
}

// GetQuestsByCreator получает квесты по создателю
func (s *QuestStorage) GetQuestsByCreator(ctx context.Context, creator string) ([]questrepo.QuestDTO, error) {
	var quests []questrepo.QuestDTO
	err := s.db.WithContext(ctx).Where("creator = ?", creator).Find(&quests).Error
	return quests, err
}

// GetQuestsByAssignee получает квесты по назначенному пользователю
func (s *QuestStorage) GetQuestsByAssignee(ctx context.Context, assignee string) ([]questrepo.QuestDTO, error) {
	var quests []questrepo.QuestDTO
	err := s.db.WithContext(ctx).Where("assignee = ?", assignee).Find(&quests).Error
	return quests, err
}

// CountQuests подсчитывает количество квестов
func (s *QuestStorage) CountQuests(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&questrepo.QuestDTO{}).Count(&count).Error
	return count, err
}

// CountQuestsByStatus подсчитывает количество квестов по статусу
func (s *QuestStorage) CountQuestsByStatus(ctx context.Context, status quest.Status) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&questrepo.QuestDTO{}).Where("status = ?", string(status)).Count(&count).Error
	return count, err
}

// DeleteQuestByID удаляет квест по ID
func (s *QuestStorage) DeleteQuestByID(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&questrepo.QuestDTO{}, "id = ?", id).Error
}

// DeleteAllQuests удаляет все квесты
func (s *QuestStorage) DeleteAllQuests(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec("DELETE FROM quests").Error
}
