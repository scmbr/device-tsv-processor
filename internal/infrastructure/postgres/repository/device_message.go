package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	dberrs "github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/errs"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/models"
)

type deviceMessageRepo struct {
	db *sqlx.DB
}

func NewDeviceMessageRepository(db *sqlx.DB) *deviceMessageRepo {
	return &deviceMessageRepo{db: db}
}

func (r *deviceMessageRepo) GetByDeviceGUID(ctx context.Context, deviceGUID string, offset, limit int) ([]*domain.DeviceMessage, int, error) {
	const op = "device_message.repo.get_by_device_guid"

	var args []any
	var sb strings.Builder

	sb.WriteString(`
        SELECT id, device_id, inv_id, msg_id, text, context, class, level, area,
               addr, block, type, bit, invert_bit, created_at
        FROM device_message
        WHERE device_id = (
            SELECT id FROM device WHERE unit_guid = $1
        )
    `)
	args = append(args, deviceGUID)

	sb.WriteString(" ORDER BY created_at ASC")
	if limit > 0 {
		args = append(args, limit)
		sb.WriteString(fmt.Sprintf(" LIMIT $%d", len(args)))
	}
	if offset > 0 {
		args = append(args, offset)
		sb.WriteString(fmt.Sprintf(" OFFSET $%d", len(args)))
	}

	var rows []*models.DeviceMessage
	if err := r.db.SelectContext(ctx, &rows, sb.String(), args...); err != nil {
		return nil, 0, dberrs.Map(err, op)
	}

	var total int
	if err := r.db.GetContext(ctx, &total,
		"SELECT COUNT(*) FROM device_message WHERE device_id = (SELECT id FROM device WHERE unit_guid = $1)",
		deviceGUID,
	); err != nil {
		return nil, 0, dberrs.Map(err, op)
	}

	out := make([]*domain.DeviceMessage, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.ToDomain())
	}

	return out, total, nil
}
func (r *deviceMessageRepo) BulkInsert(ctx context.Context, messages []*domain.DeviceMessage) error {
	const op = "device_message.repo.bulk_insert"

	if len(messages) == 0 {
		return nil
	}
	exec := GetExecFromCtx(ctx, r.db)
	valueStrings := make([]string, 0, len(messages))
	valueArgs := make([]interface{}, 0, len(messages)*14)

	for i, m := range messages {
		n := i * 14
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12, n+13, n+14,
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

	if _, err := exec.ExecContext(ctx, query, valueArgs...); err != nil {
		return dberrs.Map(err, op)
	}

	return nil
}
