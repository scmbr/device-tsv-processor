package queue

import "context"

type DocumentTask struct {
	DocumentID  int64
	UnitGUID    string
	FilePath    *string
	FileType    string
	Attempts    int
	MaxAttempts int
}

type DocumentQueue interface {
	Publish(ctx context.Context, task DocumentTask) error
	Consume(ctx context.Context) (<-chan DocumentTask, error)
	Ack(ctx context.Context, task DocumentTask) error
	Nack(ctx context.Context, task DocumentTask, requeue bool) error
}
