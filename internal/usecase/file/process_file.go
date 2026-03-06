package file

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
)

type ProcessFile struct {
	fileRepo    repository.FileRecordRepository
	messageRepo repository.DeviceMessageRepository
	errorRepo   repository.ParseErrorRepository
	txManager   repository.TxManager
	parser      TSVParser
}

func NewProcessFile(
	fileRepo repository.FileRecordRepository,
	messageRepo repository.DeviceMessageRepository,
	errorRepo repository.ParseErrorRepository,
	txManager repository.TxManager,
	parser TSVParser,
) *ProcessFile {
	return &ProcessFile{
		fileRepo:    fileRepo,
		messageRepo: messageRepo,
		errorRepo:   errorRepo,
		txManager:   txManager,
		parser:      parser,
	}
}

type ProcessFileInput struct {
	FileID int64
	Path   string
}

func (uc *ProcessFile) Execute(ctx context.Context, input ProcessFileInput) error {
	const op = "usecase.process_file"

	if err := uc.fileRepo.UpdateStatus(ctx, input.FileID, domain.FileRecordStatusProcessing); err != nil {
		return errs.Wrap(op, err)
	}

	records, parseErrors, err := uc.parser.Parse(ctx, input.Path)
	if err != nil {
		_ = uc.fileRepo.MarkFailed(ctx, input.FileID, err.Error())
		return errs.Wrap(op, err)
	}

	err = uc.txManager.WithTx(ctx, func(txCtx context.Context) error {

		if len(records) > 0 {
			if err := uc.messageRepo.BulkInsert(txCtx, records); err != nil {
				return err
			}
		}

		if len(parseErrors) > 0 {
			if err := uc.errorRepo.BulkInsert(txCtx, parseErrors); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		_ = uc.fileRepo.MarkFailed(ctx, input.FileID, err.Error())
		return errs.Wrap(op, err)
	}

	if err := uc.fileRepo.UpdateStatus(ctx, input.FileID, domain.FileRecordStatusProcessed); err != nil {
		return errs.Wrap(op, err)
	}

	return nil
}
