package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type DeviceMessageRepository interface {
	GetByDeviceGUID(ctx context.Context, deviceGUID string, offset, limit int) ([]*domain.DeviceMessage, int, error)
	BulkInsert(ctx context.Context, messages []*domain.DeviceMessage) error
}
