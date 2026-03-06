package file_processor

import (
	"context"
	"sync"

	"github.com/scmbr/device-tsv-processor/internal/queue"
	"github.com/scmbr/device-tsv-processor/internal/usecase/file"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type ProcessWorker struct {
	processUC           *file.ProcessFile
	incrementAttemptsUC *file.IncrementFileAttempts
	markErrorUC         *file.MarkFileAsError
	queue               queue.FileQueue
	maxAttempts         int
}

type TaskResult struct {
	Task  queue.FileTask
	Error error
}

func NewProcessWorker(
	processUC *file.ProcessFile,
	incrementAttemptsUC *file.IncrementFileAttempts,
	markErrorUC *file.MarkFileAsError,
	queue queue.FileQueue,
	maxAttempts int,
) *ProcessWorker {
	return &ProcessWorker{
		processUC:           processUC,
		incrementAttemptsUC: incrementAttemptsUC,
		markErrorUC:         markErrorUC,
		queue:               queue,
		maxAttempts:         maxAttempts,
	}
}

func (w *ProcessWorker) StartPool(ctx context.Context, workerCount int) error {
	logger.Info("process worker pool started", map[string]interface{}{
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
					"fileID": res.Task.FileID,
				})
				continue
			}

			w.queue.Nack(ctx, res.Task, true)
			logger.Error("file processing failed", res.Error, map[string]interface{}{
				"fileID": res.Task.FileID,
				"path":   res.Task.FullPath,
			})
			continue
		}

		w.queue.Ack(ctx, res.Task)
	}

	logger.Info("process worker pool stopped", nil)
	return nil
}

func (w *ProcessWorker) handleTask(ctx context.Context, t queue.FileTask, workerID int) TaskResult {
	select {
	case <-ctx.Done():
		return TaskResult{Task: t, Error: ctx.Err()}
	default:
	}

	updatedAttempts, err := w.incrementAttemptsUC.Execute(ctx, file.IncrementFileAttemptsInput{
		FileID: t.FileID,
	})
	if err != nil {
		return TaskResult{Task: t, Error: err}
	}

	if updatedAttempts > w.maxAttempts {
		if err := w.markErrorUC.Execute(ctx, file.MarkFileAsErrorInput{
			FileID:   t.FileID,
			Attempts: updatedAttempts,
		}); err != nil {
			logger.Error("failed to mark file as error", err, map[string]interface{}{
				"fileID": t.FileID,
			})
		}
		return TaskResult{Task: t, Error: nil}
	}

	input := file.ProcessFileInput{
		FileID: t.FileID,
		Path:   t.FullPath,
	}

	if err := w.processUC.Execute(ctx, input); err != nil {
		return TaskResult{Task: t, Error: err}
	}

	logger.Info("file processed successfully", map[string]interface{}{
		"fileID":   t.FileID,
		"path":     t.FullPath,
		"workerID": workerID,
	})

	return TaskResult{Task: t, Error: nil}
}
