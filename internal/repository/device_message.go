package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type DeviceMessageRepository interface {
	CreateBatch(ctx context.Context, messages []*domain.DeviceMessage) error
	GetByDeviceGUID(ctx context.Context, deviceGUID string, offset, limit int) ([]*domain.DeviceMessage, int, error)
	Create(ctx context.Context, message *domain.DeviceMessage) error
	BulkInsert(ctx context.Context, messages []*domain.DeviceMessage) error
}
