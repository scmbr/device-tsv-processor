package file

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/repository"
)

type GetPendingFiles struct {
	fileRepo  repository.FileRecordRepository
	batchSize int
}

func NewGetPendingFiles(fileRepo repository.FileRecordRepository, batchSize int) *GetPendingFiles {
	return &GetPendingFiles{fileRepo: fileRepo, batchSize: batchSize}
}

func (uc *GetPendingFiles) Execute(ctx context.Context) ([]*domain.FileRecord, error) {
	return uc.fileRepo.GetPending(ctx, uc.batchSize)
}
