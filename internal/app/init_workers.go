package app

import (
	"github.com/scmbr/device-tsv-processor/internal/app/config"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq/queue"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/workers/document_generator"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/workers/file_processor"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/workers/file_scanner"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/workers/task_queuer"
	"github.com/scmbr/device-tsv-processor/internal/usecase"
)

type Workers struct {
	Scan      *file_scanner.ScanWorker
	Process   *file_processor.ProcessWorker
	Queue     *task_queuer.QueueWorker
	Generator *document_generator.GeneratorWorker
}

func initWorkers(cfg *config.Config, ucs *usecase.UseCases, queues *queue.Queues) *Workers {
	return &Workers{
		Scan: file_scanner.NewScanWorker(
			ucs.ScanDirectory,
			cfg.Workers.ScanDirectoryInterval,
		),
		Process: file_processor.NewProcessWorker(
			ucs.ProcessFile,
			ucs.IncrementFileAttempts,
			ucs.MarkFileAsError,
			queues.FileQueue,
			cfg.MaxAttempts,
		),
		Queue: task_queuer.NewQueueWorker(
			ucs.EnqueueFileProcessing,
			ucs.EnqueueDocumentGenerating,
			cfg.Workers.EnqueueTasksInterval,
		),
		Generator: document_generator.NewGeneratorWorker(
			ucs.GenerateDocument,
			ucs.IncrementDocumentAttempts,
			ucs.MarkDocumentAsError,
			queues.DocumentQueue,
			cfg.MaxAttempts,
		),
	}
}
