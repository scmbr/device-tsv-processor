package file_scanner

import (
	"context"
	"time"

	"github.com/scmbr/device-tsv-processor/internal/usecase"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type ScanWorker struct {
	scanUC    *usecase.ScanDirectory
	dirPath   string
	batchSize int
	interval  time.Duration
}

func NewScanWorker(
	scanUC *usecase.ScanDirectory,
	dirPath string,
	batchSize int,
	interval time.Duration,
) *ScanWorker {
	return &ScanWorker{
		scanUC:    scanUC,
		dirPath:   dirPath,
		batchSize: batchSize,
		interval:  interval,
	}
}

func (w *ScanWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("scan worker stopped", nil)
			return
		case <-ticker.C:
			w.scanOnce(ctx)
		}
	}
}

func (w *ScanWorker) scanOnce(ctx context.Context) {
	input := usecase.ScanDirectoryInput{
		DirPath: w.dirPath,
	}

	err := w.scanUC.Execute(ctx, input)
	if err != nil {
		logger.Error("scan directory failed", err, map[string]interface{}{
			"directory": w.dirPath,
		})
		return
	}

	logger.Info("directory scanned successfully", map[string]interface{}{
		"directory": w.dirPath,
	})
}
