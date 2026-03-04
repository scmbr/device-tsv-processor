package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	dberrs "github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/errs"
)

type txManager struct {
	db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) *txManager {
	return &txManager{db: db}
}

func (t *txManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	const op = "tx_manager.with_tx"

	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return dberrs.Map(err, op)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return dberrs.Map(err, op)
	}

	return nil
}

type txKey struct{}

func GetExecFromCtx(ctx context.Context, db *sqlx.DB) sqlx.ExtContext {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok && tx != nil {
		return tx
	}
	return db
}
