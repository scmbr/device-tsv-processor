package file_processor

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/queue"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/internal/usecase"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type ProcessWorker struct {
	processUC   *usecase.ProcessFile
	fileRepo    repository.FileRecordRepository
	queue       queue.FileQueue
	maxAttempts int
}

func NewProcessWorker(
	processUC *usecase.ProcessFile,
	fileRepo repository.FileRecordRepository,
	queue queue.FileQueue,
	maxAttempts int,
	enqueueUC *usecase.EnqueueDocumentGenerating,
) *ProcessWorker {
	return &ProcessWorker{
		processUC:   processUC,
		fileRepo:    fileRepo,
		queue:       queue,
		maxAttempts: maxAttempts,
	}
}

func (w *ProcessWorker) Start(ctx context.Context) error {
	tasks, err := w.queue.Consume(ctx)
	if err != nil {
		return err
	}

	for dt := range tasks {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := w.handleTask(ctx, dt); err != nil {
			w.queue.Nack(ctx, dt, true)
			continue
		}

		w.queue.Ack(ctx, dt)
	}

	return nil
}

func (w *ProcessWorker) handleTask(ctx context.Context, t queue.FileTask) error {

	if t.Attempts >= w.maxAttempts {
		logger.Error("file reached max attempts, marking as error", nil, map[string]interface{}{
			"fileID":   t.FileID,
			"path":     t.FullPath,
			"attempts": t.Attempts,
		})

		_ = w.fileRepo.UpdateStatus(ctx, t.FileID, domain.FileRecordStatusError)
		return nil
	}

	t.Attempts++
	if err := w.fileRepo.UpdateAttempts(ctx, t.FileID, t.Attempts); err != nil {
		logger.Error("failed to update attempts", err, map[string]interface{}{
			"fileID": t.FileID,
		})
		return err
	}

	input := usecase.ProcessFileInput{
		FileID: t.FileID,
		Path:   t.FullPath,
	}

	if err := w.processUC.Execute(ctx, input); err != nil {
		logger.Error("failed to process file", err, map[string]interface{}{
			"fileID":  t.FileID,
			"path":    t.FullPath,
			"attempt": t.Attempts,
		})
		return err
	}
	logger.Info("file processed successfully", map[string]interface{}{
		"fileID":  t.FileID,
		"path":    t.FullPath,
		"attempt": t.Attempts,
	})

	return nil
}
