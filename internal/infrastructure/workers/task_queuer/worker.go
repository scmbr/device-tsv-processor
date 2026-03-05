package task_queuer

import (
	"context"
	"time"

	"github.com/scmbr/device-tsv-processor/internal/usecase"
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
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *QueueWorker) runOnce(ctx context.Context) {
	if w.fileEnqueueUC != nil {
		_ = w.fileEnqueueUC.Execute(ctx)
	}

	if w.documentEnqueueUC != nil {
		_ = w.documentEnqueueUC.Execute(ctx)
	}
}
