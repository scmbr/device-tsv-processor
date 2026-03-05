package usecase

import (
	"github.com/scmbr/device-tsv-processor/internal/repository"
)

type UseCases struct {
	EnqueueFileProcessing *EnqueueFileProcessing
	GenerateDocument      *GenerateDocument
	GetDeviceMessages     *GetDeviceMessages
	ProcessFile           *ProcessFile
	ScanDirectory         *ScanDirectory
}
type UseCaseConfig struct {
	Repos        *repository.Repositories
	BatchSize    int
	OutputDir    string
	PDFGenerator PDFGenerator
	Parser       TSVParser
}

func NewUseCases(cfg UseCaseConfig) *UseCases {
	return &UseCases{
		EnqueueFileProcessing: NewEnqueueFileProcessing(
			cfg.Repos.FileRecordRepository,
			cfg.Repos.FileQueue,
			cfg.BatchSize,
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
	}
}
