package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type DeviceMessageRepository interface {
	GetByDeviceGUID(ctx context.Context, deviceGUID string, offset, limit int) ([]*domain.DeviceMessage, int, error)
	BulkInsert(ctx context.Context, messages []*domain.DeviceMessage) error
}

type deviceMessageRepo struct {
	db *sqlx.DB
}

func NewDeviceMessageRepository(db *sqlx.DB) DeviceMessageRepository {
	return &deviceMessageRepo{db: db}
}

func (r *deviceMessageRepo) GetByDeviceGUID(ctx context.Context, deviceGUID string, offset, limit int) ([]*domain.DeviceMessage, int, error) {
	const op = "device_message.repo.get_by_device_guid"

	var deviceID int64
	if err := r.db.GetContext(ctx, &deviceID, "SELECT id FROM device WHERE unit_guid = $1", deviceGUID); err != nil {
		return nil, 0, errs.E(errs.KindNotFound, "DEVICE_NOT_FOUND", op, err.Error(), nil, nil)
	}

	var total int
	if err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM device_message WHERE device_id = $1", deviceID); err != nil {
		return nil, 0, errs.Wrap(op, err)
	}

	query := `
        SELECT id, device_id, inv_id, msg_id, text, context, class, level, area,
               addr, block, type, bit, invert_bit, created_at
        FROM device_message
        WHERE device_id = $1
        ORDER BY created_at ASC
        OFFSET $2 LIMIT $3;
    `
	rows := []*domain.DeviceMessage{}
	if err := r.db.SelectContext(ctx, &rows, query, deviceID, offset, limit); err != nil {
		return nil, 0, errs.Wrap(op, err)
	}

	return rows, total, nil
}

func (r *deviceMessageRepo) BulkInsert(ctx context.Context, messages []*domain.DeviceMessage) error {
	const op = "device_message.repo.bulk_insert"

	if len(messages) == 0 {
		return nil
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for i, m := range messages {
		n := i * 14
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12, n+13, n+14, n+15,
			),
		)
		valueArgs = append(valueArgs,
			m.DeviceID,
			m.InvID,
			m.MsgID,
			m.Text,
			m.Context,
			m.Class,
			m.Level,
			m.Area,
			m.Addr,
			m.Block,
			m.Type,
			m.Bit,
			m.InvertBit,
			m.CreatedAt,
		)
	}

	query := fmt.Sprintf(`
        INSERT INTO device_message (
            device_id, inv_id, msg_id, text, context, class, level, area,
            addr, block, type, bit, invert_bit, created_at
        ) VALUES %s
    `, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return errs.Wrap(op, err)
	}
	return nil
}
