package ports

import (
	"context"
)

// Repositories holds all repository instances for a transaction
type Repositories struct {
	Quest    QuestRepository
	Location LocationRepository
	Event    EventPublisher
}

// TransactionManager defines the interface for managing database transactions
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}
