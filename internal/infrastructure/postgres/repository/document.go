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

type documentRepo struct {
	db *sqlx.DB
}

func NewDocumentRepository(db *sqlx.DB) *documentRepo {
	return &documentRepo{db: db}
}

func (r *documentRepo) Create(ctx context.Context, doc *domain.Document) error {
	const op = "document.repo.create"

	query := `
        INSERT INTO documents (unit_guid, file_path, file_type, status, attempts, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

	now := time.Now()
	doc.CreatedAt = now
	doc.UpdatedAt = &now

	err := r.db.QueryRowContext(
		ctx,
		query,
		doc.UnitGUID,
		doc.FilePath,
		doc.FileType,
		string(doc.Status),
		doc.Attempts,
		doc.CreatedAt,
		doc.UpdatedAt,
	).Scan(&doc.ID)
	if err != nil {
		return dberrs.Map(err, op)
	}
	return nil
}

func (r *documentRepo) Exists(ctx context.Context, unitGUID string) (bool, error) {
	const op = "document.repo.exists"

	query := `
        SELECT 1
        FROM documents
        WHERE unit_guid = $1
        LIMIT 1
    `

	var dummy int
	err := r.db.GetContext(ctx, &dummy, query, unitGUID)
	if err != nil {
		if errs.IsKind(dberrs.Map(err, op), dberrs.KindNotFound) {
			return false, nil
		}
		return false, dberrs.Map(err, op)
	}

	return true, nil
}

func (r *documentRepo) GetPending(ctx context.Context, batchSize int) ([]*domain.Document, error) {
	const op = "document.repo.get_pending"

	query := `
        SELECT id, unit_guid, file_path, file_type, status, attempts, created_at, updated_at
        FROM documents
        WHERE status = $1
        ORDER BY created_at
        LIMIT $2
    `

	var models []*models.Document
	if err := r.db.SelectContext(ctx, &models, query, string(domain.DocumentStatusPending), batchSize); err != nil {
		return nil, dberrs.Map(err, op)
	}

	out := make([]*domain.Document, 0, len(models))
	for _, m := range models {
		out = append(out, m.ToDomain())
	}

	return out, nil
}

func (r *documentRepo) UpdateStatus(ctx context.Context, id int64, status domain.DocumentStatus) error {
	const op = "document.repo.update_status"

	query := `
        UPDATE documents
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
			"DOCUMENT_NOT_FOUND",
			op,
			"document not found",
			nil,
			nil,
		)
	}

	return nil
}

func (r *documentRepo) UpdateAttempts(ctx context.Context, id int64, attempts int) error {
	const op = "document.repo.update_attempts"

	query := `
        UPDATE documents
        SET attempts = $1,
            updated_at = NOW()
        WHERE id = $2
    `

	_, err := r.db.ExecContext(ctx, query, attempts, id)
	if err != nil {
		return dberrs.Map(err, op)
	}

	return nil
}
func (r *documentRepo) IncrementAttempts(ctx context.Context, documentID int64) (int, error) {
	const op = "document.repo.increment_attempts"

	query := `
        UPDATE documents
        SET attempts = attempts + 1,
            updated_at = NOW()
        WHERE id = $1
        RETURNING attempts
    `

	var updatedAttempts int
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(&updatedAttempts)
	if err != nil {
		return 0, dberrs.Map(err, op)
	}

	return updatedAttempts, nil
}
