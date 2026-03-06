package file

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type MarkFileAsError struct {
	fileRepo repository.FileRecordRepository
}

func NewMarkFileAsError(fileRepo repository.FileRecordRepository) *MarkFileAsError {
	return &MarkFileAsError{fileRepo: fileRepo}
}

type MarkFileAsErrorInput struct {
	FileID   int64
	Attempts int
}

func (uc *MarkFileAsError) Execute(ctx context.Context, input MarkFileAsErrorInput) error {
	const op = "usecase.mark_file_as_error"

	logger.Error("file reached max attempts, marking as error", nil, map[string]interface{}{
		"fileID":   input.FileID,
		"attempts": input.Attempts,
	})

	if err := uc.fileRepo.UpdateStatus(ctx, input.FileID, domain.FileRecordStatusError); err != nil {
		return errs.Wrap(op, err)
	}
	return nil
}
