package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *domain.Document) error
	Exists(ctx context.Context, unitGUID string) (bool, error)
	UpdateStatus(ctx context.Context, id int64, status domain.DocumentStatus) error
	GetPending(ctx context.Context, batchSize int) ([]*domain.Document, error)
	UpdateAttempts(ctx context.Context, id int64, attempts int) error
}
