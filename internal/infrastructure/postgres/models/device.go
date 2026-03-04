package models

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type Device struct {
	ID          int64      `db:"id,primarykey,autoincrement"`
	UnitGUID    string     `db:"unit_guid,unique,notnull"`
	InvID       string     `db:"inv_id,notnull"`
	MQTT        string     `db:"mqtt,omitempty"`
	Status      string     `db:"status,notnull"`
	ProcessedAt *time.Time `db:"processed_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at,notnull"`
}

func (m *Device) ToDomain() *domain.Device {
	if m == nil {
		return nil
	}

	return &domain.Device{
		ID:          m.ID,
		GUID:        m.UnitGUID,
		InvID:       m.InvID,
		MQTT:        m.MQTT,
		Status:      m.Status,
		ProcessedAt: m.ProcessedAt,
		CreatedAt:   m.CreatedAt,
	}
}
