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

// UnitOfWorkFactory creates a fresh UnitOfWork and an EventPublisher bound to the
// same transactional tracker. Callers are responsible for managing the
// transaction lifecycle (Begin/Commit/Rollback) for the returned UnitOfWork.
type UnitOfWorkFactory func() (UnitOfWork, EventPublisher, error)
