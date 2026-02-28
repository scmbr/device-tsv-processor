package usecase

import (
	"context"
	"os"
	"strings"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
)

type FileRecordUseCase struct {
	unitRepo       repository.UnitRepository
	messageRepo    repository.UnitMessageRepository
	fileRepo       repository.FileRecordRepository
	parseErrorRepo repository.ParseErrorRepository
	outboxRepo     repository.OutboxRepository
}

func NewFileRecordUseCase(
	unitRepo repository.UnitRepository,
	messageRepo repository.UnitMessageRepository,
	fileRepo repository.FileRecordRepository,
	parseErrorRepo repository.ParseErrorRepository,
	outboxRepo repository.OutboxRepository,
) *FileRecordUseCase {
	return &FileRecordUseCase{
		unitRepo:       unitRepo,
		messageRepo:    messageRepo,
		fileRepo:       fileRepo,
		parseErrorRepo: parseErrorRepo,
		outboxRepo:     outboxRepo,
	}
}

type ScanDirectoryInput struct {
	dirPath   string
	batchSize int
}

func (u *FileRecordUseCase) ScanDirectory(ctx context.Context, in ScanDirectoryInput) error {
	const op = "file_record.usecase.scan_directory"
	var err error
	files, err := os.ReadDir(in.dirPath)
	if err != nil {
		return errs.Wrap(op, err)
	}
	records := make([]*domain.FileRecord, 0, len(files))
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".tsv") {
			continue
		}
		fileDomain, err := domain.NewFileRecord(f.Name(), domain.FileRecordStatusPending)
		if err != nil {
			return errs.Wrap(op, err)
		}
		records = append(records, fileDomain)
	}
	if len(records) == 0 {
		return nil
	}
	for start := 0; start < len(records); start += in.batchSize {
		end := start + in.batchSize
		if end > len(records) {
			end = len(records)
		}
		chunk := records[start:end]
		err = u.fileRepo.BatchInsert(ctx, chunk)
		if err != nil {
			return errs.Wrap(op, err)
		}
	}
	return nil
}
