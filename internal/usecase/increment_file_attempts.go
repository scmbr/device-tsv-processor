package usecase

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type IncrementFileAttempts struct {
	fileRepo repository.FileRecordRepository
}

func NewIncrementFileAttempts(fileRepo repository.FileRecordRepository) *IncrementFileAttempts {
	return &IncrementFileAttempts{fileRepo: fileRepo}
}

type IncrementFileAttemptsInput struct {
	FileID int64
}

func (uc *IncrementFileAttempts) Execute(ctx context.Context, input IncrementFileAttemptsInput) (int, error) {
	const op = "usecase.increment_file_attempts"

	attempts, err := uc.fileRepo.IncrementAttempts(ctx, input.FileID)
	if err != nil {
		logger.Error("failed to update attempts", err, map[string]interface{}{
			"fileID": input.FileID,
		})
		return 0, errs.Wrap(op, err)
	}

	return attempts, nil
}
