package document

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type IncrementDocumentAttempts struct {
	documentRepo repository.DocumentRepository
}

func NewIncrementDocumentAttempts(documentRepo repository.DocumentRepository) *IncrementDocumentAttempts {
	return &IncrementDocumentAttempts{documentRepo: documentRepo}
}

type IncrementDocumentAttemptsInput struct {
	DocumentID int64
}

func (uc *IncrementDocumentAttempts) Execute(ctx context.Context, input IncrementDocumentAttemptsInput) (int, error) {
	const op = "usecase.increment_document_attempts"

	attempts, err := uc.documentRepo.IncrementAttempts(ctx, input.DocumentID)
	if err != nil {
		logger.Error("failed to update attempts", err, map[string]interface{}{
			"documentID": input.DocumentID,
		})
		return 0, errs.Wrap(op, err)
	}

	return attempts, nil
}
