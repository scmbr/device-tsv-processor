package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type FileRecordRepository interface {
	Create(ctx context.Context, file *domain.FileRecord) error
	Exists(ctx context.Context, filename string) (bool, error)
	BatchInsert(ctx context.Context, chunk []*domain.FileRecord) error
	UpdateStatus(ctx context.Context, id int, status domain.FileRecordStatus) error
	MarkFailed(ctx context.Context, id int, error string) error
	GetPending(ctx context.Context, batchSize int) ([]*domain.FileRecord, error)
}

type fileRecordRepo struct {
	db *sqlx.DB
}

func NewFileRecordRepository(db *sqlx.DB) FileRecordRepository {
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

	return r.db.QueryRowContext(
		ctx,
		query,
		file.Filename,
		file.Status,
		file.ErrorMessage,
		file.CreatedAt,
		file.UpdatedAt,
	).Scan(&file.ID)
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
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errs.Wrap(op, err)
	}

	return true, nil
}
func (r *fileRecordRepo) BatchInsert(ctx context.Context, chunk []*domain.FileRecord) error {
	const op = "file_record.repo.batch_insert"

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
			file.Status,
			file.ErrorMessage,
			file.CreatedAt,
			file.UpdatedAt,
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
func (r *fileRecordRepo) UpdateStatus(ctx context.Context, id int, status domain.FileRecordStatus) error {
	const op = "file_record.repo.update_status"

	query := `
        UPDATE file_records
        SET status = $1,
            updated_at = $2
        WHERE id = $3
    `

	res, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return errs.Wrap(op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errs.Wrap(op, err)
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
		domain.FileRecordStatusError,
		errorMsg,
		time.Now(),
		id,
	)
	if err != nil {
		return errs.Wrap(op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errs.Wrap(op, err)
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

	var files []*domain.FileRecord
	if err := r.db.SelectContext(
		ctx,
		&files,
		query,
		domain.FileRecordStatusPending,
		batchSize,
	); err != nil {
		return nil, errs.Wrap(op, err)
	}

	return files, nil
}
