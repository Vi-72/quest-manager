package storage

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quest-manager/internal/adapters/out/postgres/questrepo"
	"quest-manager/internal/core/domain/model/quest"
)

// QuestStorage provides methods for working with quests in tests
type QuestStorage struct {
	db *gorm.DB
}

// NewQuestStorage creates new QuestStorage
func NewQuestStorage(db *gorm.DB) *QuestStorage {
	return &QuestStorage{db: db}
}

// GetQuestByID gets quest from database by ID
func (s *QuestStorage) GetQuestByID(ctx context.Context, id uuid.UUID) (*questrepo.QuestDTO, error) {
	var questDTO questrepo.QuestDTO
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&questDTO).Error
	if err != nil {
		return nil, err
	}
	return &questDTO, nil
}

// GetAllQuests gets all quests from database
func (s *QuestStorage) GetAllQuests(ctx context.Context) ([]questrepo.QuestDTO, error) {
	var quests []questrepo.QuestDTO
	err := s.db.WithContext(ctx).Find(&quests).Error
	return quests, err
}

// GetQuestsByStatus gets quests by status
func (s *QuestStorage) GetQuestsByStatus(ctx context.Context, status quest.Status) ([]questrepo.QuestDTO, error) {
	var quests []questrepo.QuestDTO
	err := s.db.WithContext(ctx).Where("status = ?", string(status)).Find(&quests).Error
	return quests, err
}

// GetQuestsByCreator gets quests by creator
func (s *QuestStorage) GetQuestsByCreator(ctx context.Context, creator string) ([]questrepo.QuestDTO, error) {
	var quests []questrepo.QuestDTO
	err := s.db.WithContext(ctx).Where("creator = ?", creator).Find(&quests).Error
	return quests, err
}

// GetQuestsByAssignee gets quests by assigned user
func (s *QuestStorage) GetQuestsByAssignee(ctx context.Context, assignee string) ([]questrepo.QuestDTO, error) {
	var quests []questrepo.QuestDTO
	err := s.db.WithContext(ctx).Where("assignee = ?", assignee).Find(&quests).Error
	return quests, err
}

// CountQuests counts number of quests
func (s *QuestStorage) CountQuests(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&questrepo.QuestDTO{}).Count(&count).Error
	return count, err
}

// CountQuestsByStatus counts number of quests by status
func (s *QuestStorage) CountQuestsByStatus(ctx context.Context, status quest.Status) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&questrepo.QuestDTO{}).Where("status = ?", string(status)).Count(&count).Error
	return count, err
}

// DeleteQuestByID deletes quest by ID
func (s *QuestStorage) DeleteQuestByID(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&questrepo.QuestDTO{}, "id = ?", id).Error
}

// DeleteAllQuests deletes all quests
func (s *QuestStorage) DeleteAllQuests(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec("DELETE FROM quests").Error
}
