package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type UnitMessageRepository interface {
	CreateBatch(ctx context.Context, messages []*domain.UnitMessage) error
	GetByUnitGUID(ctx context.Context, unitGUID string, offset, limit int) ([]*domain.UnitMessage, error)
	Create(ctx context.Context, message *domain.UnitMessage) error
}
