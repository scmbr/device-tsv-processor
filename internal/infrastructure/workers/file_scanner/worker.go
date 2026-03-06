package file_scanner

import (
	"context"
	"time"

	"github.com/scmbr/device-tsv-processor/internal/usecase/file"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type ScanWorker struct {
	scanUC *file.ScanDirectory

	interval time.Duration
}

func NewScanWorker(
	scanUC *file.ScanDirectory,
	interval time.Duration,
) *ScanWorker {
	return &ScanWorker{
		scanUC:   scanUC,
		interval: interval,
	}
}

func (w *ScanWorker) Start(ctx context.Context) error {
	logger.Info("scan worker started", map[string]interface{}{
		"interval_ms": w.interval,
	})
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	w.scanOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			logger.Info("scan worker stopped", nil)
			return ctx.Err()
		case <-ticker.C:
			w.scanOnce(ctx)
		}
	}
}

func (w *ScanWorker) scanOnce(ctx context.Context) {
	scanCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	start := time.Now()
	if err := w.scanUC.Execute(scanCtx); err != nil {
		logger.Error("scan directory failed", err, map[string]interface{}{
			"duration_ms": time.Since(start).Milliseconds(),
		})
		return
	}
	logger.Info("scan directory completed", map[string]interface{}{
		"duration_ms": time.Since(start).Milliseconds(),
	})
}
