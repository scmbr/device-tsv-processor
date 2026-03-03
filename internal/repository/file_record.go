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
	UpdateStatus(ctx context.Context, id int, status domain.FileRecordStatus) error
	MarkFailed(ctx context.Context, id int, error string) error
	GetPending(ctx context.Context, batchSize int) ([]*domain.FileRecord, error)
}
type FileQueue interface {
	Enqueue(fileID int) error
	Dequeue() (fileID int, err error)
}
