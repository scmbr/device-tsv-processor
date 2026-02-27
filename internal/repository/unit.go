package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type UnitRepository interface {
	GetByGUID(ctx context.Context, guid string) (*domain.Unit, error)
	Create(ctx context.Context, unit *domain.Unit) error
	Update(ctx context.Context, unit *domain.Unit) error
	List(ctx context.Context, offset, limit int) ([]*domain.Unit, error)
}
