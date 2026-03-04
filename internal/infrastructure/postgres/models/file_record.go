package models

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type FileRecord struct {
	ID           int64      `db:"id,primarykey,autoincrement"`
	Filename     string     `db:"filename,unique,notnull"`
	CreatedAt    time.Time  `db:"created_at,notnull"`
	ProcessedAt  *time.Time `db:"processed_at"`
	Status       string     `db:"status,notnull"`
	ErrorMessage *string    `db:"error_message,omitempty"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

func (m *FileRecord) ToDomain() *domain.FileRecord {
	if m == nil {
		return nil
	}

	return &domain.FileRecord{
		ID:           m.ID,
		Filename:     m.Filename,
		CreatedAt:    m.CreatedAt,
		ProcessedAt:  m.ProcessedAt,
		Status:       domain.FileRecordStatus(m.Status),
		ErrorMessage: m.ErrorMessage,
		UpdatedAt:    m.UpdatedAt,
	}
}
