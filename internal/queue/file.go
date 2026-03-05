package queue

import "context"

type FileTask struct {
	FileID      int64
	FullPath    string
	Filename    string
	Attempts    int
	MaxAttempts int
}

type FileQueue interface {
	Publish(ctx context.Context, task FileTask) error
	Consume(ctx context.Context) (<-chan FileTask, error)
	Ack(ctx context.Context, task FileTask) error
	Nack(ctx context.Context, task FileTask, requeue bool) error
}
