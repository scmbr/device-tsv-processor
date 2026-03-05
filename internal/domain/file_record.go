package domain

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type FileRecordStatus string

var (
	FileRecordStatusProcessed  FileRecordStatus = "PROCESSED"
	FileRecordStatusQueued     FileRecordStatus = "QUEUED"
	FileRecordStatusProcessing FileRecordStatus = "PROCESSING"
	FileRecordStatusPending    FileRecordStatus = "PENDING"
	FileRecordStatusError      FileRecordStatus = "ERROR"
)

type FileRecord struct {
	ID           int64
	Filename     string
	FullPath     string
	ProcessedAt  *time.Time
	CreatedAt    time.Time
	Status       FileRecordStatus
	ErrorMessage *string
	UpdatedAt    *time.Time
	Attempts     int
}

func NewFileRecord(fileName, fullPath string, status FileRecordStatus) (*FileRecord, error) {
	fileRecord := &FileRecord{
		Filename: fileName,
		FullPath: fullPath,
		Status:   status,
	}
	if err := fileRecord.Validate(); err != nil {
		return nil, err
	}
	return fileRecord, nil
}
func (f *FileRecord) Validate() error {
	const op = "file_record.entity.validate"
	fields := map[string]string{}
	if f.Filename == "" {
		fields["filename"] = "is required"
	}
	if f.FullPath == "" {
		fields["full_path"] = "is required"
	}
	if f.Status == "" {
		fields["status"] = "is required"
	}
	if len(fields) > 0 {
		return errs.E(errs.KindInvalid, "FILE_RECORD_INVALID", op, "invalid file record", fields, nil)
	}
	return nil
}
