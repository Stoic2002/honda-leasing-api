package database

import (
	"context"

	"honda-leasing-api/internal/domain/contract"

	"gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "gorm_tx"

type uow struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) contract.UnitOfWork {
	return &uow{db: db}
}

// Do executes fn within a database transaction.
// The transaction is injected into the context so that repositories
// can retrieve it via GetTxFromContext.
func (u *uow) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey, tx)
		return fn(txCtx)
	})
}

// GetTxFromContext retrieves the GORM transaction from context.
// If no transaction is found, it returns the provided fallback db.
func GetTxFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return fallback
}
