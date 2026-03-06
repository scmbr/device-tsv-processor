package file

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type ProcessFile struct {
	fileRepo     repository.FileRecordRepository
	deviceRepo   repository.DeviceRepository
	messageRepo  repository.DeviceMessageRepository
	documentRepo repository.DocumentRepository
	errorRepo    repository.ParseErrorRepository
	txManager    repository.TxManager
	parser       TSVParser
}

func NewProcessFile(
	fileRepo repository.FileRecordRepository,
	messageRepo repository.DeviceMessageRepository,
	deviceRepo repository.DeviceRepository,
	errorRepo repository.ParseErrorRepository,
	documentRepo repository.DocumentRepository,
	txManager repository.TxManager,
	parser TSVParser,
) *ProcessFile {
	return &ProcessFile{
		fileRepo:     fileRepo,
		messageRepo:  messageRepo,
		errorRepo:    errorRepo,
		documentRepo: documentRepo,
		deviceRepo:   deviceRepo,
		txManager:    txManager,
		parser:       parser,
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
	if len(records) == 0 && len(parseErrors) == 0 {
		return uc.fileRepo.MarkFailed(ctx, input.FileID, "file is empty or invalid")
	}
	for _, m := range records {
		device, err := uc.deviceRepo.CreateIfNotExists(ctx, &domain.Device{
			GUID:   m.UnitGUID,
			InvID:  m.InvID,
			MQTT:   m.MQTT,
			Status: "",
		})
		if err != nil {
			_ = uc.fileRepo.MarkFailed(ctx, input.FileID, err.Error())
			return errs.Wrap(op, err)
		}
		m.DeviceID = device.ID
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
	for _, m := range records {
		doc, err := domain.NewDocument(
			m.UnitGUID,
			"pdf",
			domain.DocumentStatusPending,
			nil,
		)
		if err != nil {

			logger.Error("failed to create document object", err, map[string]interface{}{
				"unitGUID": m.UnitGUID,
			})
			continue
		}

		_, err = uc.documentRepo.CreateIfNotExists(ctx, doc)
		if err != nil {

			logger.Error("failed to create document in DB", err, map[string]interface{}{
				"unitGUID": m.UnitGUID,
			})
			continue
		}
	}
	if err := uc.fileRepo.UpdateStatus(ctx, input.FileID, domain.FileRecordStatusProcessed); err != nil {
		return errs.Wrap(op, err)
	}

	return nil
}
