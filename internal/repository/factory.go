package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/repository"
)

type Repositories struct {
	DeviceMessageRepository DeviceMessageRepository
	DeviceRepository        DeviceRepository
	FileRecordRepository    FileRecordRepository
	FileQueue               FileQueue
	ParseErrorRepository    ParseErrorRepository
	TxManager               TxManager
}

func NewRepositories(db *sqlx.DB, fileQueue FileQueue) *Repositories {
	return &Repositories{
		DeviceMessageRepository: repository.NewDeviceMessageRepository(db),
		DeviceRepository:        repository.NewDeviceRepository(db),
		FileRecordRepository:    repository.NewFileRecordRepository(db),
		FileQueue:               fileQueue,
		ParseErrorRepository:    repository.NewParseErrorRepository(db),
		TxManager:               repository.NewTxManager(db),
	}
}
