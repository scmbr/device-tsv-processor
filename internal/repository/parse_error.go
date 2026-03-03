package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type ParseErrorRepository interface {
	BulkInsert(ctx context.Context, errors []*domain.ParseError) error
}
