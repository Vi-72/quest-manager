package ports

import (
	"context"
)

type UnitOfWork interface {
	Begin(ctx context.Context)
	Commit(ctx context.Context) error
	Rollback() error
	QuestRepository() QuestRepository
	LocationRepository() LocationRepository
}
