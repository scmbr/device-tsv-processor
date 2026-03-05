package usecase

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type TSVParser interface {
	Parse(ctx context.Context, path string) ([]*domain.DeviceMessage, []*domain.ParseError, error)
}
