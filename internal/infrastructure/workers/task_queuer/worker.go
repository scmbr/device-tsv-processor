package task_queuer

import (
	"context"
	"time"

	"github.com/scmbr/device-tsv-processor/internal/usecase"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type QueueWorker struct {
	fileEnqueueUC     *usecase.EnqueueFileProcessing
	documentEnqueueUC *usecase.EnqueueDocumentGenerating
	interval          time.Duration
}

func NewQueueWorker(
	fileUC *usecase.EnqueueFileProcessing,
	documentUC *usecase.EnqueueDocumentGenerating,
	interval time.Duration,
) *QueueWorker {
	return &QueueWorker{
		fileEnqueueUC:     fileUC,
		documentEnqueueUC: documentUC,
		interval:          interval,
	}
}

func (w *QueueWorker) Start(ctx context.Context) error {
	logger.Info("queue worker started", map[string]interface{}{
		"interval": w.interval,
	})
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.runOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			logger.Info("queue worker stopped", nil)
			return ctx.Err()

		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *QueueWorker) runOnce(ctx context.Context) {
	fileCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

	defer cancel()
	if err := w.fileEnqueueUC.Execute(fileCtx); err != nil {
		logger.Error("enqueue file failed", err, nil)
	}

	docCtx, cancelDoc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelDoc()
	if err := w.documentEnqueueUC.Execute(docCtx); err != nil {
		logger.Error("enqueue document failed", err, nil)
	}
}
