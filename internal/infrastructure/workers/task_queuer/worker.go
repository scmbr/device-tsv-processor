package task_queuer

import (
	"context"
	"time"

	"github.com/scmbr/device-tsv-processor/internal/usecase/document"
	"github.com/scmbr/device-tsv-processor/internal/usecase/file"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type QueueWorker struct {
	fileEnqueueUC     *file.EnqueueFileProcessing
	documentEnqueueUC *document.EnqueueDocumentGenerating
	interval          time.Duration
}

func NewQueueWorker(
	fileUC *file.EnqueueFileProcessing,
	documentUC *document.EnqueueDocumentGenerating,
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
		"interval": w.interval.String(),
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
	fileCtx, cancelFile := context.WithTimeout(ctx, 30*time.Second)

	defer cancelFile()
	if err := w.fileEnqueueUC.Execute(fileCtx); err != nil {
		logger.Error("enqueue file failed", err, nil)
	}

	docCtx, cancelDoc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelDoc()
	if err := w.documentEnqueueUC.Execute(docCtx); err != nil {
		logger.Error("enqueue document failed", err, nil)
	}
}
