package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type FileRecordRepository interface {
	Create(ctx context.Context, file *domain.FileRecord) error
	Exists(ctx context.Context, filename string) (bool, error)
	List(ctx context.Context, offset, limit int) ([]*domain.FileRecord, error)
	BatchInsert(ctx context.Context, chunk []*domain.FileRecord) error
	ClaimPendingBatch(ctx context.Context, batchSize int) ([]*domain.FileRecord, error)
	MarkProcessedBatch(ctx context.Context, ids []int) error
}
