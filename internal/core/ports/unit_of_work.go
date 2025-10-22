package ports

import (
	"context"
)

type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback() error
	QuestRepository() QuestRepository
	LocationRepository() LocationRepository
}

// UnitOfWorkFactory creates new UnitOfWork instances for each request
// This ensures thread safety by avoiding shared state between concurrent requests
type UnitOfWorkFactory interface {
	CreateUnitOfWork() (UnitOfWork, error)
}
