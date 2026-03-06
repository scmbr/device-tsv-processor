package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type DeviceRepository interface {
	CreateIfNotExists(ctx context.Context, device *domain.Device) (*domain.Device, error)
}
