package domain

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type DocumentStatus string

var (
	DocumentStatusGenerated  DocumentStatus = "GENERATED"
	DocumentStatusQueued     DocumentStatus = "QUEUED"
	DocumentStatusGenerating DocumentStatus = "GENERATING"
	DocumentStatusPending    DocumentStatus = "PENDING"
	DocumentStatusError      DocumentStatus = "ERROR"
)

type Document struct {
	ID        int64
	UnitGUID  string
	FilePath  *string
	FileType  string
	Status    DocumentStatus
	Attempts  int
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewDocument(unitGUID string, fileType string, status DocumentStatus, filePath *string) (*Document, error) {
	now := time.Now()
	doc := &Document{
		UnitGUID:  unitGUID,
		FileType:  fileType,
		Status:    status,
		FilePath:  filePath,
		CreatedAt: now,
		UpdatedAt: nil,
		Attempts:  0,
	}
	if err := doc.Validate(); err != nil {
		return nil, err
	}
	return doc, nil
}

func (d *Document) Validate() error {
	const op = "document.entity.validate"
	fields := map[string]string{}

	if d.UnitGUID == "" {
		fields["unit_guid"] = "is required"
	}
	if d.FileType == "" {
		fields["file_type"] = "is required"
	}
	if d.Status == "" {
		fields["status"] = "is required"
	}

	if len(fields) > 0 {
		return errs.E(errs.KindInvalid, "DOCUMENT_INVALID", op, "invalid document", fields, nil)
	}
	return nil
}
