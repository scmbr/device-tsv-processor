package usecase

import (
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/repository"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq/queue"
	"github.com/scmbr/device-tsv-processor/internal/usecase/device_message"
	"github.com/scmbr/device-tsv-processor/internal/usecase/document"
	"github.com/scmbr/device-tsv-processor/internal/usecase/file"
)

type UseCases struct {
	Document      *DocumentUseCases
	File          *FileUseCases
	DeviceMessage *DeviceMessageUseCases
}
type DocumentUseCases struct {
	EnqueueDocumentGenerating *document.EnqueueDocumentGenerating
	GenerateDocument          *document.GenerateDocument
	MarkDocumentAsError       *document.MarkDocumentAsError
	IncrementDocumentAttempts *document.IncrementDocumentAttempts
}
type FileUseCases struct {
	EnqueueFileProcessing *file.EnqueueFileProcessing
	ProcessFile           *file.ProcessFile
	ScanDirectory         *file.ScanDirectory

	MarkFileAsError       *file.MarkFileAsError
	IncrementFileAttempts *file.IncrementFileAttempts
}
type DeviceMessageUseCases struct {
	GetDeviceMessages *device_message.GetDeviceMessages
}
type UseCaseConfig struct {
	Repos        *repository.Repositories
	Queues       *queue.Queues
	BatchSize    int
	OutputDir    string
	InputDir     string
	PDFGenerator document.PDFGenerator
	Parser       file.TSVParser
	MaxAttempts  int
}

func NewUseCases(cfg UseCaseConfig) *UseCases {
	return &UseCases{
		Document: &DocumentUseCases{
			GenerateDocument: document.NewGenerateDocument(
				cfg.Repos.DeviceMessageRepository,
				cfg.Repos.ParseErrorRepository,
				cfg.Repos.DocumentRepository,
				cfg.OutputDir,
				cfg.PDFGenerator,
			),
			EnqueueDocumentGenerating: document.NewEnqueueDocumentProcessing(
				cfg.Repos.DocumentRepository,
				cfg.Queues.DocumentQueue,
				cfg.BatchSize,
				cfg.MaxAttempts,
			),

			MarkDocumentAsError: document.NewMarkDocumentAsError(
				cfg.Repos.DocumentRepository,
			),
			IncrementDocumentAttempts: document.NewIncrementDocumentAttempts(
				cfg.Repos.DocumentRepository,
			),
		},
		File: &FileUseCases{
			EnqueueFileProcessing: file.NewEnqueueFileProcessing(
				cfg.Repos.FileRecordRepository,
				cfg.Queues.FileQueue,
				cfg.BatchSize,
				cfg.MaxAttempts,
			),
			ProcessFile: file.NewProcessFile(
				cfg.Repos.FileRecordRepository,
				cfg.Repos.DeviceMessageRepository,
				cfg.Repos.DeviceRepository,
				cfg.Repos.ParseErrorRepository,
				cfg.Repos.DocumentRepository,
				cfg.Repos.TxManager,
				cfg.Parser,
			),
			ScanDirectory: file.NewScanDirectory(
				cfg.Repos.FileRecordRepository,
				cfg.InputDir,
				cfg.BatchSize,
			),
			MarkFileAsError: file.NewMarkFileAsError(
				cfg.Repos.FileRecordRepository,
			),
			IncrementFileAttempts: file.NewIncrementFileAttempts(
				cfg.Repos.FileRecordRepository,
			),
		},
		DeviceMessage: &DeviceMessageUseCases{
			GetDeviceMessages: device_message.NewGetDeviceMessages(cfg.Repos.DeviceMessageRepository),
		},
	}
}
