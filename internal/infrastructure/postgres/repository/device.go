package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	dberrs "github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/errs"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/models"
)

type deviceRepo struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) *deviceRepo {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) GetByUnitGUID(ctx context.Context, unitGUID string) (*domain.Device, error) {
	const op = "device.repo.get_by_unit_guid"

	var m models.Device
	query := `
        SELECT id, unit_guid, inv_id, mqtt, status, processed_at, created_at
        FROM devices
        WHERE unit_guid = $1
    `
	if err := r.db.GetContext(ctx, &m, query, unitGUID); err != nil {
		return nil, dberrs.Map(err, op)
	}

	return m.ToDomain(), nil
}

func (r *deviceRepo) Create(ctx context.Context, device *domain.Device) (*domain.Device, error) {
	const op = "device.repo.create"

	query := `
        INSERT INTO devices (unit_guid, inv_id, mqtt, status, processed_at, created_at)
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
		device.Status,
		device.ProcessedAt,
		device.CreatedAt,
	).Scan(&id)
	if err != nil {
		return nil, dberrs.Map(err, op)
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

	return nil, dberrs.Map(err, op)
}
