package queue

import (
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq"
	"github.com/scmbr/device-tsv-processor/internal/queue"
)

type Queues struct {
	FileQueue     queue.FileQueue
	DocumentQueue queue.DocumentQueue
}
type QueuesConfig struct {
	Client            *rabbitmq.Client
	FileQueueName     string
	DocumentQueueNAme string
}

func NewQueues(cfg QueuesConfig) (*Queues, error) {
	fileQueue, err := NewRabbitMQFileQueue(cfg.Client, cfg.FileQueueName)

	if err != nil {
		return nil, err
	}
	documentQueue, err := NewRabbitMQDocumentQueue(cfg.Client, cfg.DocumentQueueNAme)
	return &Queues{
		FileQueue:     fileQueue,
		DocumentQueue: documentQueue,
	}, nil
}
