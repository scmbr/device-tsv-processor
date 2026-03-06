package file

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/queue"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type EnqueueFileProcessing struct {
	fileRepo    repository.FileRecordRepository
	queue       queue.FileQueue
	batchSize   int
	maxAttempts int
}

func NewEnqueueFileProcessing(fileRepo repository.FileRecordRepository, queue queue.FileQueue, batchSize int, maxAttempts int) *EnqueueFileProcessing {
	return &EnqueueFileProcessing{
		fileRepo:    fileRepo,
		queue:       queue,
		batchSize:   batchSize,
		maxAttempts: maxAttempts,
	}
}

func (uc *EnqueueFileProcessing) Execute(ctx context.Context) error {
	const op = "usecase.enqueue_file_processing"

	files, err := uc.fileRepo.GetPending(ctx, uc.batchSize)
	if err != nil {
		return errs.Wrap(op, err)
	}

	for _, f := range files {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		task := queue.FileTask{
			FileID:      f.ID,
			FullPath:    f.FullPath,
			Filename:    f.Filename,
			Attempts:    0,
			MaxAttempts: uc.maxAttempts,
		}

		if err := uc.queue.Publish(ctx, task); err != nil {
			logger.Error("Failed to enqueue task", err, map[string]interface{}{"file_id": f.ID, "filename": f.Filename})
			continue
		}

		if err := uc.fileRepo.UpdateStatus(ctx, f.ID, domain.FileRecordStatusQueued); err != nil {
			logger.Error("Failed to update file status to Queued", err, map[string]interface{}{"file_id": f.ID, "filename": f.Filename})
			continue
		}
	}

	return nil
}
