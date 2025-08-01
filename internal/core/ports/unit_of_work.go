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
