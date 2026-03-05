package usecase

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
)

type GenerateDocument struct {
	messageRepo  repository.DeviceMessageRepository
	errorRepo    repository.ParseErrorRepository
	outputDir    string
	docGenerator PDFGenerator
}

func NewGenerateDocument(
	messageRepo repository.DeviceMessageRepository,
	errorRepo repository.ParseErrorRepository,
	outputDir string,
	docGenerator PDFGenerator,
) *GenerateDocument {
	return &GenerateDocument{
		messageRepo:  messageRepo,
		errorRepo:    errorRepo,
		outputDir:    outputDir,
		docGenerator: docGenerator,
	}
}

type GenerateDocumentInput struct {
	UnitGUID string
}

func (uc *GenerateDocument) Execute(ctx context.Context, input GenerateDocumentInput) error {
	const op = "usecase.generate_document"

	messages, _, err := uc.messageRepo.GetByDeviceGUID(ctx, input.UnitGUID, 0, 0) // 0,0 = все записи
	if err != nil {
		return errs.Wrap(op, err)
	}

	if len(messages) == 0 {
		return errs.Wrap(op, err)
	}

	outputPath := filepath.Join(uc.outputDir, fmt.Sprintf("%s.pdf", input.UnitGUID))
	if err := uc.docGenerator.Generate(outputPath, messages); err != nil {
		if pe, e := domain.NewParseError(fmt.Sprintf("doc_%s", input.UnitGUID), 0, err.Error()); e == nil {
			_ = uc.errorRepo.BulkInsert(ctx, []*domain.ParseError{pe})
		}
		return errs.Wrap(op, err)
	}

	return nil
}
