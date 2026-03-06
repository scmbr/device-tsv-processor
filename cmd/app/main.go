package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/scmbr/device-tsv-processor/internal/app"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

const configsDir = "configs"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	application, err := app.New("configs")
	if err != nil {
		logger.Error("app init failed", err, nil)
		os.Exit(1)
	}

	if err := application.Run(ctx); err != nil {
		logger.Error("app stopped with error", err, nil)
		os.Exit(1)
	}

	logger.Info("app stopped", nil)
}
