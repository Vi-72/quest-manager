package mocks

import (
	"context"
	"fmt"
	"sync"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"

	"github.com/google/uuid"
)

// MockQuestRepository is an in-memory implementation of QuestRepository for contract testing
type MockQuestRepository struct {
	quests map[uuid.UUID]quest.Quest
	mu     sync.RWMutex
}

func NewMockQuestRepository() *MockQuestRepository {
	return &MockQuestRepository{
		quests: make(map[uuid.UUID]quest.Quest),
	}
}

func (m *MockQuestRepository) GetByID(ctx context.Context, questID uuid.UUID) (quest.Quest, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	q, exists := m.quests[questID]
	if !exists {
		return quest.Quest{}, fmt.Errorf("quest with id %s not found", questID.String())
	}
	return q, nil
}

func (m *MockQuestRepository) Save(ctx context.Context, q quest.Quest) error {
	_ = ctx // unused in mock
	m.mu.Lock()
	defer m.mu.Unlock()
	m.quests[q.ID()] = q
	return nil
}

func (m *MockQuestRepository) FindAll(ctx context.Context) ([]quest.Quest, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []quest.Quest
	for _, q := range m.quests {
		result = append(result, q)
	}
	return result, nil
}

func (m *MockQuestRepository) FindByStatus(ctx context.Context, status quest.Status) ([]quest.Quest, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []quest.Quest
	for _, q := range m.quests {
		if q.Status == status {
			result = append(result, q)
		}
	}
	return result, nil
}

func (m *MockQuestRepository) FindByBoundingBox(ctx context.Context, bbox kernel.BoundingBox) ([]quest.Quest, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []quest.Quest
	for _, q := range m.quests {
		// Check if either target or execution location is within bounding box
		if m.isWithinBoundingBox(q.TargetLocation, bbox) || m.isWithinBoundingBox(q.ExecutionLocation, bbox) {
			result = append(result, q)
		}
	}
	return result, nil
}

func (m *MockQuestRepository) FindByAssignee(ctx context.Context, userID string) ([]quest.Quest, error) {
	_ = ctx // unused in mock
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []quest.Quest
	for _, q := range m.quests {
		if q.Assignee != nil && *q.Assignee == userID {
			result = append(result, q)
		}
	}
	return result, nil
}

func (m *MockQuestRepository) isWithinBoundingBox(coord kernel.GeoCoordinate, bbox kernel.BoundingBox) bool {
	return coord.Lat >= bbox.MinLat && coord.Lat <= bbox.MaxLat &&
		coord.Lon >= bbox.MinLon && coord.Lon <= bbox.MaxLon
}

// Helper methods for testing
func (m *MockQuestRepository) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.quests = make(map[uuid.UUID]quest.Quest)
}

func (m *MockQuestRepository) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.quests)
}
