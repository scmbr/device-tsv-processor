package document

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type MarkDocumentAsError struct {
	documentRepo repository.DocumentRepository
}

func NewMarkDocumentAsError(documentRepo repository.DocumentRepository) *MarkDocumentAsError {
	return &MarkDocumentAsError{documentRepo: documentRepo}
}

type MarkDocumentAsErrorInput struct {
	DocumentID int64

	Attempts int
}

func (uc *MarkDocumentAsError) Execute(ctx context.Context, input MarkDocumentAsErrorInput) error {
	const op = "usecase.mark_document_as_error"

	logger.Error("document reached max attempts, marking as error", nil, map[string]interface{}{
		"documentID": input.DocumentID,
		"attempts":   input.Attempts,
	})

	if err := uc.documentRepo.UpdateStatus(ctx, input.DocumentID, domain.DocumentStatusError); err != nil {
		return errs.Wrap(op, err)
	}
	return nil
}
