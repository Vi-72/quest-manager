package postgres

import (
	"context"

	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/adapters/out/postgres/locationrepo"
	"quest-manager/internal/adapters/out/postgres/questrepo"
	"quest-manager/internal/core/ports"

	"gorm.io/gorm"
)

// TransactionManager coordinates multi-aggregate transactions
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// RunInTransaction executes a function within a database transaction
// All repository operations within the function will use the same transaction
func (tm *TransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context, repos ports.Repositories) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create repository instances that use the transaction
		repos := ports.Repositories{
			Quest:    questrepo.NewRepository(tx),
			Location: locationrepo.NewRepository(tx),
			Event:    eventrepo.NewRepository(tx),
		}
		return fn(ctx, repos)
	})
}
