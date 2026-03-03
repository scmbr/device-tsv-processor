package models

import "time"

type FileRecord struct {
	ID           int        `db:"id,primarykey,autoincrement"`
	Filename     string     `db:"filename,unique,notnull"`
	CreatedAt    time.Time  `db:"created_at,notnull"`
	ProcessedAt  *time.Time `db:"processed_at"`
	Status       string     `db:"status,notnull"`
	ErrorMessage *string    `db:"error_message,omitempty"`
	UpdatedAt    *time.Time `db:"updated_at"`
}
