package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type DeviceRepository interface {
	GetByGUID(ctx context.Context, guid string) (*domain.Device, error)
	Create(ctx context.Context, device *domain.Device) error
	CreateIfNotExists(ctx context.Context, device *domain.Device) (*domain.Device, error)
}
