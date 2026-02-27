package repository

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/models"
)

type OutboxRepository interface {
	Enqueue(ctx context.Context, event *models.OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]*models.OutboxEvent, error)
	MarkAsProcessed(ctx context.Context, eventID string) error
}
