package models

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type Document struct {
	ID        int64      `db:"id,primarykey,autoincrement"`
	UnitGUID  string     `db:"unit_guid,notnull"`
	FilePath  *string    `db:"file_path,omitempty"`
	FileType  string     `db:"file_type,notnull"`
	Status    string     `db:"status,notnull"`
	Attempts  int        `db:"attempts,notnull"`
	CreatedAt time.Time  `db:"created_at,notnull"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (m *Document) ToDomain() *domain.Document {
	if m == nil {
		return nil
	}

	return &domain.Document{
		ID:        m.ID,
		UnitGUID:  m.UnitGUID,
		FilePath:  m.FilePath,
		FileType:  m.FileType,
		Status:    domain.DocumentStatus(m.Status),
		Attempts:  m.Attempts,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
