package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type ProcessedFileRepository interface {
	Create(ctx context.Context, file *domain.ProcessedFile) error
	Exists(ctx context.Context, filename string) (bool, error)
	List(ctx context.Context, offset, limit int) ([]*domain.ProcessedFile, error)
}
