package postgres

import (
	"context"
	"quest-manager/internal/adapters/out/postgres/questrepo"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"

	"gorm.io/gorm"
)

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	tx              *gorm.DB
	db              *gorm.DB
	questRepository ports.QuestRepository
}

func NewUnitOfWork(db *gorm.DB) (ports.UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	uow := &UnitOfWork{db: db}

	questRepo, err := questrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.questRepository = questRepo

	return uow, nil
}

func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}

func (u *UnitOfWork) Db() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

func (u *UnitOfWork) Begin(ctx context.Context) {
	u.tx = u.db.WithContext(ctx).Begin()
}

func (u *UnitOfWork) Rollback() error {
	if u.tx != nil {
		err := u.tx.Rollback().Error
		u.tx = nil
		return err
	}
	return nil
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}
	u.tx = nil
	return nil
}

// Геттеры репозиториев
func (u *UnitOfWork) QuestRepository() ports.QuestRepository {
	return u.questRepository
}
