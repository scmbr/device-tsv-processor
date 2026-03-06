package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type FileRecordRepository interface {
	Create(ctx context.Context, file *domain.FileRecord) error
	Exists(ctx context.Context, filename string) (bool, error)
	BatchInsert(ctx context.Context, chunk []*domain.FileRecord) error
	UpdateStatus(ctx context.Context, id int64, status domain.FileRecordStatus) error
	MarkFailed(ctx context.Context, id int64, error string) error
	GetPending(ctx context.Context, batchSize int) ([]*domain.FileRecord, error)
	IncrementAttempts(ctx context.Context, id int64) (int, error)
}
