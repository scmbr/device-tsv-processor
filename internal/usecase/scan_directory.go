package usecase

import (
	"context"
	"os"
	"strings"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type ScanDirectory struct {
	fileRepo repository.FileRecordRepository
}

func NewScanDirectory(fileRepo repository.FileRecordRepository) *ScanDirectory {
	return &ScanDirectory{fileRepo: fileRepo}
}

type ScanDirectoryInput struct {
	DirPath   string
	BatchSize int
}

func (uc *ScanDirectory) Execute(ctx context.Context, input ScanDirectoryInput) error {
	const op = "usecase.scan_directory"

	dirEntries, err := os.ReadDir(input.DirPath)
	if err != nil {
		return errs.Wrap(op, err)
	}

	if len(dirEntries) == 0 {
		logger.Info("no files found in directory", map[string]interface{}{
			"directory_path": input.DirPath,
		})
		return nil
	}

	batch := make([]*domain.FileRecord, 0, input.BatchSize)
	insertBatch := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := uc.fileRepo.BatchInsert(ctx, batch); err != nil {
			return errs.Wrap(op, err)
		}
		batch = batch[:0]
		return nil
	}

	for _, entry := range dirEntries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tsv") {
			continue
		}

		exists, err := uc.fileRepo.Exists(ctx, entry.Name())
		if err != nil {
			return errs.Wrap(op, err)
		}
		if exists {
			continue
		}

		fileDomain, err := domain.NewFileRecord(entry.Name(), domain.FileRecordStatusPending)
		if err != nil {
			logger.Error("failed to create FileRecord", err, map[string]interface{}{
				"filename": entry.Name(),
			})
			continue
		}

		batch = append(batch, fileDomain)
		if len(batch) >= input.BatchSize {
			if err := insertBatch(); err != nil {
				return err
			}
		}
	}

	if err := insertBatch(); err != nil {
		return err
	}

	return nil
}
