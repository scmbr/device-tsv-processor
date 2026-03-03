package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type DeviceRepository interface {
	GetByUnitGUID(ctx context.Context, unitGUID string) (*domain.Device, error)
	Create(ctx context.Context, device *domain.Device) (*domain.Device, error)
	CreateIfNotExists(ctx context.Context, device *domain.Device) (*domain.Device, error)
}

type deviceRepo struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) DeviceRepository {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) GetByUnitGUID(ctx context.Context, unitGUID string) (*domain.Device, error) {
	const op = "device.repo.get_by_unit_guid"

	device := &domain.Device{}
	query := `
        SELECT id, guid, inv_id, mqtt, processed_at, status, created_at
        FROM device
        WHERE guid = $1
    `
	if err := r.db.GetContext(ctx, device, query, unitGUID); err != nil {
		return nil, errs.E(errs.KindNotFound, "DEVICE_NOT_FOUND", op, err.Error(), nil, nil)
	}

	return device, nil
}

func (r *deviceRepo) Create(ctx context.Context, device *domain.Device) (*domain.Device, error) {
	const op = "device.repo.create"

	query := `
        INSERT INTO device (guid, inv_id, mqtt, processed_at, status, created_at)
        VALUES ($1,$2,$3,$4,$5,$6)
        RETURNING id
    `
	now := time.Now()
	device.CreatedAt = now

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		device.GUID,
		device.InvID,
		device.MQTT,
		device.ProcessedAt,
		device.Status,
		device.CreatedAt,
	).Scan(&id)
	if err != nil {
		return nil, errs.Wrap(op, err)
	}
	device.ID = id
	return device, nil
}

func (r *deviceRepo) CreateIfNotExists(ctx context.Context, device *domain.Device) (*domain.Device, error) {
	const op = "device.repo.create_if_not_exists"

	existing, err := r.GetByUnitGUID(ctx, device.GUID)
	if err == nil {
		return existing, nil
	}

	if errs.IsKind(err, errs.KindNotFound) {
		return r.Create(ctx, device)
	}

	return nil, errs.Wrap(op, err)
}
