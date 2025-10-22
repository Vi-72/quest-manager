package tests

import (
	"gorm.io/gorm"

	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/core/ports"
)

// TestUnitOfWorkFactory creates UnitOfWork instances for integration tests
type TestUnitOfWorkFactory struct {
	db *gorm.DB
}

func NewTestUnitOfWorkFactory(db *gorm.DB) *TestUnitOfWorkFactory {
	return &TestUnitOfWorkFactory{db: db}
}

func (f *TestUnitOfWorkFactory) CreateUnitOfWork() (ports.UnitOfWork, error) {
	return postgres.NewUnitOfWork(f.db)
}
