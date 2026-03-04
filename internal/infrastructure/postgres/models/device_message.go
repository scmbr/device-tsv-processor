package models

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type DeviceMessage struct {
	ID        int64     `db:"id,primarykey,autoincrement"`
	DeviceID  int64     `db:"device_id,notnull"`
	InvID     string    `db:"inv_id,notnull"`
	MsgID     string    `db:"msg_id,unique,notnull"`
	Text      string    `db:"text"`
	Context   string    `db:"context"`
	Class     string    `db:"class"`
	Level     int       `db:"level"`
	Area      string    `db:"area"`
	Addr      string    `db:"addr"`
	Block     string    `db:"block"`
	Type      string    `db:"type"`
	Bit       int       `db:"bit"`
	InvertBit bool      `db:"invert_bit"`
	CreatedAt time.Time `db:"created_at,notnull"`
}

func (m *DeviceMessage) ToDomain() *domain.DeviceMessage {
	if m == nil {
		return nil
	}

	return &domain.DeviceMessage{
		ID:        m.ID,
		DeviceID:  m.DeviceID,
		InvID:     m.InvID,
		MsgID:     m.MsgID,
		Text:      m.Text,
		Context:   m.Context,
		Class:     m.Class,
		Level:     m.Level,
		Area:      m.Area,
		Addr:      m.Addr,
		Block:     m.Block,
		Type:      m.Type,
		Bit:       m.Bit,
		InvertBit: m.InvertBit,
		CreatedAt: m.CreatedAt,
	}
}
