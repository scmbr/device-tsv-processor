package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/repository"
)

type Repositories struct {
	DeviceMessageRepository repository.DeviceMessageRepository
	DeviceRepository        repository.DeviceRepository
	FileRecordRepository    repository.FileRecordRepository
	ParseErrorRepository    repository.ParseErrorRepository
	DocumentRepository      repository.DocumentRepository
	TxManager               repository.TxManager
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		DeviceMessageRepository: NewDeviceMessageRepository(db),
		DeviceRepository:        NewDeviceRepository(db),
		FileRecordRepository:    NewFileRecordRepository(db),
		ParseErrorRepository:    NewParseErrorRepository(db),
		DocumentRepository:      NewDocumentRepository(db),
		TxManager:               NewTxManager(db),
	}
}
