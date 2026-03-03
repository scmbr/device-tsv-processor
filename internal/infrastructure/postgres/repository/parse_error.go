package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type ParseErrorRepository interface {
	BulkInsert(ctx context.Context, errors []*domain.ParseError) error
}

type parseErrorRepo struct {
	db *sqlx.DB
}

func NewParseErrorRepository(db *sqlx.DB) ParseErrorRepository {
	return &parseErrorRepo{db: db}
}

func (r *parseErrorRepo) BulkInsert(ctx context.Context, errors []*domain.ParseError) error {
	const op = "parse_error.repo.bulk_insert"

	if len(errors) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errs.Wrap(op, err)
	}
	defer tx.Rollback()

	query := `
        INSERT INTO parse_errors (filename,file_id, line, raw, error, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	now := time.Now()

	for _, e := range errors {
		if e.CreatedAt.IsZero() {
			e.CreatedAt = now
		}

		_, err := tx.ExecContext(
			ctx,
			query,
			e.Filename,
			e.FileID,
			e.LineNumber,
			e.Message,
			e.CreatedAt,
		)
		if err != nil {
			return errs.Wrap(op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return errs.Wrap(op, err)
	}

	return nil
}
