package models

import "time"

type DeviceMessage struct {
	ID        int64     `db:"id,primarykey,autoincrement"`
	DeviceID  int       `db:"device_id,notnull"`
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
