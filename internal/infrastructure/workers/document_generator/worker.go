package document_generator

import (
	"context"
	"sync"

	"github.com/scmbr/device-tsv-processor/internal/queue"
	"github.com/scmbr/device-tsv-processor/internal/usecase"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type GeneratorWorker struct {
	generateUC          *usecase.GenerateDocument
	incrementAttemptsUC *usecase.IncrementDocumentAttempts
	markErrorUC         *usecase.MarkDocumentAsError
	queue               queue.DocumentQueue
	maxAttempts         int
}

type TaskResult struct {
	Task  queue.DocumentTask
	Error error
}

func NewGeneratorWorker(
	generateUC *usecase.GenerateDocument,
	incrementAttemptsUC *usecase.IncrementDocumentAttempts,
	markErrorUC *usecase.MarkDocumentAsError,
	queue queue.DocumentQueue,
	maxAttempts int,
) *GeneratorWorker {
	return &GeneratorWorker{
		generateUC:          generateUC,
		incrementAttemptsUC: incrementAttemptsUC,
		markErrorUC:         markErrorUC,
		queue:               queue,
		maxAttempts:         maxAttempts,
	}
}

func (w *GeneratorWorker) StartPool(ctx context.Context, workerCount int) error {
	logger.Info("generate worker pool started", map[string]interface{}{
		"workers": workerCount,
	})

	tasks, err := w.queue.Consume(ctx)
	if err != nil {
		return err
	}

	results := make(chan TaskResult)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-tasks:
					if !ok {
						return
					}
					res := w.handleTask(ctx, task, workerID)
					results <- res
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Error != nil {
			if res.Error == context.Canceled {
				logger.Info("task canceled due to context", map[string]interface{}{
					"unitGUID": res.Task.UnitGUID,
				})
				continue
			}

			w.queue.Nack(ctx, res.Task, true)
			logger.Error("task processing failed", res.Error, map[string]interface{}{
				"unitGUID": res.Task.UnitGUID,
				"filepath": res.Task.FilePath,
			})
			continue
		}

		w.queue.Ack(ctx, res.Task)
	}

	logger.Info("generate worker pool stopped", nil)
	return nil
}

func (w *GeneratorWorker) handleTask(ctx context.Context, t queue.DocumentTask, workerID int) TaskResult {
	select {
	case <-ctx.Done():
		return TaskResult{Task: t, Error: ctx.Err()}
	default:
	}

	updatedAttempts, err := w.incrementAttemptsUC.Execute(ctx, usecase.IncrementDocumentAttemptsInput{
		DocumentID: t.DocumentID,
	})
	if err != nil {
		return TaskResult{Task: t, Error: err}
	}

	if updatedAttempts > w.maxAttempts {
		if err := w.markErrorUC.Execute(ctx, usecase.MarkDocumentAsErrorInput{
			DocumentID: t.DocumentID,
			Attempts:   updatedAttempts,
		}); err != nil {
			logger.Error("failed to mark document as error", err, map[string]interface{}{
				"filepath": t.FilePath,
			})
		}
		return TaskResult{Task: t, Error: nil}
	}

	input := usecase.GenerateDocumentInput{
		UnitGUID: t.UnitGUID,
	}

	if err := w.generateUC.Execute(ctx, input); err != nil {
		return TaskResult{Task: t, Error: err}
	}

	logger.Info("document generated successfully", map[string]interface{}{
		"unitGUID": t.UnitGUID,
		"workerID": workerID,
	})

	return TaskResult{Task: t, Error: nil}
}
