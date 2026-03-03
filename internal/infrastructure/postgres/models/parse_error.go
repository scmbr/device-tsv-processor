package models

import "time"

type ParseErrorModel struct {
	ID         int       `db:"id,primarykey,autoincrement"`
	Filename   string    `db:"filename,notnull"`
	FileID     int       `db:"file_id,notnull"`
	LineNumber int       `db:"line_number,notnull"`
	Message    string    `db:"message,notnull"`
	CreatedAt  time.Time `db:"created_at,notnull"`
}
