package usecase

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
)

type EnqueueFileProcessing struct {
	fileRepo  repository.FileRecordRepository
	queue     repository.FileQueue
	batchSize int
}

func NewEnqueueFileProcessing(fileRepo repository.FileRecordRepository, queue repository.FileQueue, batchSize int) *EnqueueFileProcessing {
	return &EnqueueFileProcessing{
		fileRepo:  fileRepo,
		queue:     queue,
		batchSize: batchSize,
	}
}

func (uc *EnqueueFileProcessing) Execute(ctx context.Context) error {
	const op = "usecase.enqueue_file_processing"

	files, err := uc.fileRepo.GetPending(ctx, uc.batchSize)
	if err != nil {
		return errs.Wrap(op, err)
	}

	for _, f := range files {
		if err := uc.queue.Enqueue(f.ID); err != nil {
			continue
		}

		if err := uc.fileRepo.UpdateStatus(ctx, f.ID, domain.FileRecordStatusQueued); err != nil {
			continue
		}
	}

	return nil
}
