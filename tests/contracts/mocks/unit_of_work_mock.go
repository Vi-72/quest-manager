package mocks

import (
	"context"
	"fmt"

	"quest-manager/internal/core/ports"
)

// MockUnitOfWork is an in-memory implementation of UnitOfWork for contract testing
type MockUnitOfWork struct {
	questRepo    ports.QuestRepository
	locationRepo ports.LocationRepository
	inTx         bool
	shouldFail   bool
}

func NewMockUnitOfWork() *MockUnitOfWork {
	return &MockUnitOfWork{
		questRepo:    NewMockQuestRepository(),
		locationRepo: NewMockLocationRepository(),
		inTx:         false,
		shouldFail:   false,
	}
}

func (m *MockUnitOfWork) Begin(ctx context.Context) error {
	_ = ctx // unused in mock
	if m.shouldFail {
		return fmt.Errorf("mock begin error")
	}
	if m.inTx {
		return fmt.Errorf("transaction already in progress")
	}
	m.inTx = true
	return nil
}

func (m *MockUnitOfWork) Commit(ctx context.Context) error {
	_ = ctx // unused in mock
	if m.shouldFail {
		return fmt.Errorf("mock commit error")
	}
	if !m.inTx {
		return fmt.Errorf("no transaction to commit")
	}
	m.inTx = false
	return nil
}

func (m *MockUnitOfWork) Rollback() error {
	if m.shouldFail {
		return fmt.Errorf("mock rollback error")
	}
	if !m.inTx {
		return fmt.Errorf("no transaction to rollback")
	}
	m.inTx = false
	return nil
}

func (m *MockUnitOfWork) QuestRepository() ports.QuestRepository {
	return m.questRepo
}

func (m *MockUnitOfWork) LocationRepository() ports.LocationRepository {
	return m.locationRepo
}

// Helper methods for testing
func (m *MockUnitOfWork) SetShouldFail(shouldFail bool) {
	m.shouldFail = shouldFail
}

func (m *MockUnitOfWork) IsInTransaction() bool {
	return m.inTx
}

func (m *MockUnitOfWork) ClearRepositories() {
	if mockQuestRepo, ok := m.questRepo.(*MockQuestRepository); ok {
		mockQuestRepo.Clear()
	}
	if mockLocationRepo, ok := m.locationRepo.(*MockLocationRepository); ok {
		mockLocationRepo.Clear()
	}
}
