package domain

import "time"

type OutboxEvent struct {
	ID          string
	UnitGUID    string
	Aggregate   string
	Payload     []byte
	CreatedAt   time.Time
	ProcessedAt *time.Time
	Attempts    uint
	Status      string
	ErrorMsg    *string `db:"error_msg"`
}
