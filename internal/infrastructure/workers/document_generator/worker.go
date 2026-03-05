package document_generator

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/queue"
	"github.com/scmbr/device-tsv-processor/internal/usecase"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type GeneratorWorker struct {
	generateUC  *usecase.GenerateDocument
	queue       queue.DocumentQueue
	maxAttempts int
}

func NewGeneratorWorker(
	generateUC *usecase.GenerateDocument,
	queue queue.DocumentQueue,
	maxAttempts int,
) *GeneratorWorker {
	return &GeneratorWorker{
		generateUC:  generateUC,
		queue:       queue,
		maxAttempts: maxAttempts,
	}
}

func (w *GeneratorWorker) Start(ctx context.Context) error {
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

func (w *GeneratorWorker) handleTask(ctx context.Context, t queue.DocumentTask) error {

	if t.Attempts >= w.maxAttempts {
		logger.Error("document generation reached max attempts", nil, map[string]interface{}{
			"unitGUID": t.UnitGUID,
			"attempts": t.Attempts,
		})
		return nil
	}

	input := usecase.GenerateDocumentInput{
		UnitGUID: t.UnitGUID,
	}

	if err := w.generateUC.Execute(ctx, input); err != nil {
		logger.Error("failed to generate document", err, map[string]interface{}{
			"unitGUID": t.UnitGUID,
			"attempt":  t.Attempts,
		})
		return err
	}

	logger.Info("document generated successfully", map[string]interface{}{
		"unitGUID": t.UnitGUID,
	})

	return nil
}
