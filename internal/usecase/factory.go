package usecase

import (
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/repository"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq/queue"
)

type UseCases struct {
	EnqueueFileProcessing     *EnqueueFileProcessing
	GenerateDocument          *GenerateDocument
	GetDeviceMessages         *GetDeviceMessages
	ProcessFile               *ProcessFile
	ScanDirectory             *ScanDirectory
	EnqueueDocumentProcessing *EnqueueDocumentGenerating
}
type UseCaseConfig struct {
	Repos        *repository.Repositories
	Queues       *queue.Queues
	BatchSize    int
	OutputDir    string
	PDFGenerator PDFGenerator
	Parser       TSVParser
	MaxAttempts  int
}

func NewUseCases(cfg UseCaseConfig) *UseCases {
	return &UseCases{
		EnqueueFileProcessing: NewEnqueueFileProcessing(
			cfg.Repos.FileRecordRepository,
			cfg.Queues.FileQueue,
			cfg.BatchSize,
			cfg.MaxAttempts,
		),
		GenerateDocument: NewGenerateDocument(
			cfg.Repos.DeviceMessageRepository,
			cfg.Repos.ParseErrorRepository,
			cfg.OutputDir,
			cfg.PDFGenerator,
		),
		GetDeviceMessages: NewGetDeviceMessages(cfg.Repos.DeviceMessageRepository),
		ProcessFile: NewProcessFile(
			cfg.Repos.FileRecordRepository,
			cfg.Repos.DeviceMessageRepository,
			cfg.Repos.ParseErrorRepository,
			cfg.Repos.TxManager,
			cfg.Parser,
		),
		ScanDirectory: NewScanDirectory(
			cfg.Repos.FileRecordRepository,
			cfg.BatchSize,
		),
		EnqueueDocumentProcessing: NewEnqueueDocumentProcessing(
			cfg.Repos.DocumentRepository,
			cfg.Queues.DocumentQueue,
			cfg.BatchSize,
			cfg.MaxAttempts,
		),
	}
}
