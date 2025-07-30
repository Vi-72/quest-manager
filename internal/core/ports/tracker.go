package ports

import (
	"context"
	"gorm.io/gorm"
)

type Tracker interface {
	Tx() *gorm.DB
	Db() *gorm.DB
	InTx() bool
	Begin(ctx context.Context)
	Commit(ctx context.Context) error
	Rollback() error
}
