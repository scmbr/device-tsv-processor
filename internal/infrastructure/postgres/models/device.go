package models

import "time"

type Device struct {
	ID        int        `db:"id,primarykey,autoincrement"`
	UnitGUID  string     `db:"unit_guid,unique,notnull"`
	InvID     string     `db:"inv_id,notnull"`
	MQTT      string     `db:"mqtt,omitempty"`
	Status    string     `db:"status,notnull"`
	Processed *time.Time `db:"processed_at,omitempty"`
	CreatedAt time.Time  `db:"created_at,notnull"`
}
