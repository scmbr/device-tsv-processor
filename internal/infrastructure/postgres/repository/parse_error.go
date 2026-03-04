package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	dberrs "github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/errs"
)

type parseErrorRepo struct {
	db *sqlx.DB
}

func NewParseErrorRepository(db *sqlx.DB) *parseErrorRepo {
	return &parseErrorRepo{db: db}
}

func (r *parseErrorRepo) BulkInsert(ctx context.Context, errors []*domain.ParseError) error {
	const op = "parse_error.repo.bulk_insert"

	if len(errors) == 0 {
		return nil
	}

	exec := GetExecFromCtx(ctx, r.db)

	query := `
        INSERT INTO parse_errors (filename,file_id, line, raw, error, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	now := time.Now()

	for _, e := range errors {
		if e.CreatedAt.IsZero() {
			e.CreatedAt = now
		}

		_, err := exec.ExecContext(
			ctx,
			query,
			e.Filename,
			e.FileID,
			e.LineNumber,
			e.Message,
			e.CreatedAt,
		)
		if err != nil {
			return dberrs.Map(err, op)
		}
	}

	return nil
}
