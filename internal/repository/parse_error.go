package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type ParseErrorRepository interface {
	Create(ctx context.Context, err *domain.ParseError) error
	ListByFile(ctx context.Context, filename string) ([]*domain.ParseError, error)
}
