package models

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type ParseErrorModel struct {
	ID         int64     `db:"id,primarykey,autoincrement"`
	Filename   string    `db:"filename,notnull"`
	FileID     int64     `db:"file_id,notnull"`
	LineNumber int       `db:"line_number,notnull"`
	Message    string    `db:"message,notnull"`
	CreatedAt  time.Time `db:"created_at,notnull"`
}

func (m *ParseErrorModel) ToDomain() *domain.ParseError {
	if m == nil {
		return nil
	}

	return &domain.ParseError{
		ID:         m.ID,
		Filename:   m.Filename,
		FileID:     m.FileID,
		LineNumber: m.LineNumber,
		Message:    m.Message,
		CreatedAt:  m.CreatedAt,
	}
}
