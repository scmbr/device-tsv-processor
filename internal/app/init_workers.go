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
			ucs.File.ScanDirectory,
			cfg.Workers.ScanDirectoryInterval,
		),
		Process: file_processor.NewProcessWorker(
			ucs.File.ProcessFile,
			ucs.File.IncrementFileAttempts,
			ucs.File.MarkFileAsError,
			queues.FileQueue,
			cfg.MaxAttempts,
		),
		Queue: task_queuer.NewQueueWorker(
			ucs.File.EnqueueFileProcessing,
			ucs.Document.EnqueueDocumentGenerating,
			cfg.Workers.EnqueueTasksInterval,
		),
		Generator: document_generator.NewGeneratorWorker(
			ucs.Document.GenerateDocument,
			ucs.Document.IncrementDocumentAttempts,
			ucs.Document.MarkDocumentAsError,
			queues.DocumentQueue,
			cfg.MaxAttempts,
		),
	}
}
