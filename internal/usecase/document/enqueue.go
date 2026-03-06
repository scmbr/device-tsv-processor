package document

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/queue"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type EnqueueDocumentGenerating struct {
	docRepo     repository.DocumentRepository
	queue       queue.DocumentQueue
	batchSize   int
	maxAttempts int
}

func NewEnqueueDocumentProcessing(docRepo repository.DocumentRepository, queue queue.DocumentQueue, batchSize int, maxAttempts int) *EnqueueDocumentGenerating {
	return &EnqueueDocumentGenerating{
		docRepo:     docRepo,
		queue:       queue,
		batchSize:   batchSize,
		maxAttempts: maxAttempts,
	}
}

func (uc *EnqueueDocumentGenerating) Execute(ctx context.Context) error {
	const op = "usecase.enqueue_document_processing"

	docs, err := uc.docRepo.GetPending(ctx, uc.batchSize)
	if err != nil {
		return errs.Wrap(op, err)
	}

	for _, d := range docs {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		task := queue.DocumentTask{
			DocumentID:  d.ID,
			UnitGUID:    d.UnitGUID,
			FileType:    d.FileType,
			Attempts:    0,
			MaxAttempts: uc.maxAttempts,
		}

		if err := uc.queue.Publish(ctx, task); err != nil {
			logger.Error("Failed to enqueue document task", err, map[string]interface{}{"document_id": d.ID, "unit_guid": d.UnitGUID})
			continue
		}

		if err := uc.docRepo.UpdateStatus(ctx, d.ID, domain.DocumentStatusQueued); err != nil {
			logger.Error("Failed to update document status to Queued", err, map[string]interface{}{"document_id": d.ID, "unit_guid": d.UnitGUID})
			continue
		}
	}

	return nil
}
