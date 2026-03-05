package app

import (
	"github.com/scmbr/device-tsv-processor/internal/app/config"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/parser"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/pdf_document"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/repository"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq/queue"
	"github.com/scmbr/device-tsv-processor/internal/usecase"
)

func initUseCases(cfg *config.Config, repos *repository.Repositories, queues *queue.Queues) *usecase.UseCases {
	pdfGenerator := pdf_document.NewPDFGenerator()
	tsvParser := parser.NewTSVParser()
	return usecase.NewUseCases(usecase.UseCaseConfig{
		Repos:        repos,
		Queues:       queues,
		BatchSize:    cfg.BatchSize,
		OutputDir:    cfg.OutputDir,
		PDFGenerator: pdfGenerator,
		Parser:       tsvParser,
		MaxAttempts:  cfg.MaxAttempts,
	})
}
