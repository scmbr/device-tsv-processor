package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	dberrs "github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/errs"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/models"
)

type fileRecordRepo struct {
	db *sqlx.DB
}

func NewFileRecordRepository(db *sqlx.DB) *fileRecordRepo {
	return &fileRecordRepo{db: db}
}

func (r *fileRecordRepo) Create(ctx context.Context, file *domain.FileRecord) error {
	const op = "file_record.repo.create"

	query := `
        INSERT INTO file_records (filename, status, error, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	now := time.Now()
	file.CreatedAt = now
	file.UpdatedAt = &now

	err := r.db.QueryRowContext(
		ctx,
		query,
		file.Filename,
		string(file.Status),
		file.ErrorMessage,
		file.CreatedAt,
		file.UpdatedAt,
	).Scan(&file.ID)

	if err != nil {
		return dberrs.Map(err, op)
	}
	return nil
}
func (r *fileRecordRepo) Exists(ctx context.Context, filename string) (bool, error) {
	const op = "file_record.repo.exists"

	query := `
        SELECT 1
        FROM file_records
        WHERE filename = $1
        LIMIT 1
    `

	var dummy int
	err := r.db.GetContext(ctx, &dummy, query, filename)
	if err != nil {
		if errs.IsKind(dberrs.Map(err, op), dberrs.KindNotFound) {
			return false, nil
		}
		return false, dberrs.Map(err, op)
	}

	return true, nil
}
func (r *fileRecordRepo) BatchInsert(ctx context.Context, chunk []*domain.FileRecord) error {
	const op = "file_record.repo.batch_insert"
	if len(chunk) == 0 {
		return nil
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errs.Wrap(op, err)
	}
	defer tx.Rollback()

	query := `
        INSERT INTO file_records (filename, status, error, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	now := time.Now()

	for _, file := range chunk {
		file.CreatedAt = now
		file.UpdatedAt = &now

		_, err := tx.ExecContext(
			ctx,
			query,
			file.Filename,
			string(file.Status),
			file.ErrorMessage,
			file.CreatedAt,
			file.UpdatedAt,
		)
		if err != nil {
			return dberrs.Map(err, op)
		}
	}

	if err := tx.Commit(); err != nil {
		return dberrs.Map(err, op)
	}

	return nil
}
func (r *fileRecordRepo) UpdateStatus(ctx context.Context, id int, status domain.FileRecordStatus) error {
	const op = "file_record.repo.update_status"

	query := `
        UPDATE file_records
        SET status = $1,
            updated_at = $2
        WHERE id = $3
    `

	res, err := r.db.ExecContext(ctx, query, string(status), time.Now(), id)
	if err != nil {
		return dberrs.Map(err, op)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return dberrs.Map(err, op)
	}

	if affected == 0 {
		return errs.E(
			errs.KindNotFound,
			"FILE_NOT_FOUND",
			op,
			"file record not found",
			nil,
			nil,
		)
	}

	return nil
}
func (r *fileRecordRepo) MarkFailed(ctx context.Context, id int, errorMsg string) error {
	const op = "file_record.repo.mark_failed"

	query := `
        UPDATE file_records
        SET status = $1,
            error = $2,
            updated_at = $3
        WHERE id = $4
    `

	res, err := r.db.ExecContext(
		ctx,
		query,
		string(domain.FileRecordStatusError),
		errorMsg,
		time.Now(),
		id,
	)
	if err != nil {
		return dberrs.Map(err, op)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return dberrs.Map(err, op)
	}

	if affected == 0 {
		return errs.E(
			errs.KindNotFound,
			"FILE_NOT_FOUND",
			op,
			"file record not found",
			nil,
			nil,
		)
	}

	return nil
}
func (r *fileRecordRepo) GetPending(ctx context.Context, batchSize int) ([]*domain.FileRecord, error) {
	const op = "file_record.repo.get_pending"

	query := `
        SELECT id, filename, status, error, created_at, updated_at
        FROM file_records
        WHERE status = $1
        ORDER BY created_at
        LIMIT $2
    `

	var models []*models.FileRecord
	if err := r.db.SelectContext(ctx, &models, query, string(domain.FileRecordStatusPending), batchSize); err != nil {
		return nil, dberrs.Map(err, op)
	}

	out := make([]*domain.FileRecord, 0, len(models))
	for _, m := range models {
		out = append(out, m.ToDomain())
	}

	return out, nil
}
