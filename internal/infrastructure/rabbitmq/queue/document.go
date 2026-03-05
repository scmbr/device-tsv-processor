package queue

import (
	"context"
	"encoding/json"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq"
	"github.com/scmbr/device-tsv-processor/internal/queue"

	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type RabbitMQDocumentQueue struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	queue       amqp.Queue
	mu          sync.Mutex
	deliveryMap map[int64]*amqp.Delivery
}

func NewRabbitMQDocumentQueue(client *rabbitmq.Client, queueName string) (*RabbitMQDocumentQueue, error) {
	ch, err := client.Conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQDocumentQueue{
		conn:        client.Conn,
		channel:     ch,
		queue:       q,
		deliveryMap: make(map[int64]*amqp.Delivery),
	}, nil
}

func (r *RabbitMQDocumentQueue) Publish(ctx context.Context, task queue.DocumentTask) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQDocumentQueue) Consume(ctx context.Context) (<-chan queue.DocumentTask, error) {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		false, // autoAck = false
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	out := make(chan queue.DocumentTask)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var t queue.DocumentTask
				if err := json.Unmarshal(msg.Body, &t); err != nil {
					logger.Error("failed to unmarshal document task", err, map[string]interface{}{"body": string(msg.Body)})
					_ = msg.Nack(false, false)
					continue
				}

				r.mu.Lock()
				r.deliveryMap[t.DocumentID] = &msg
				r.mu.Unlock()

				out <- t
			}
		}
	}()

	return out, nil
}

func (r *RabbitMQDocumentQueue) Ack(ctx context.Context, task queue.DocumentTask) error {
	r.mu.Lock()
	delivery, ok := r.deliveryMap[task.DocumentID]
	if ok {
		delete(r.deliveryMap, task.DocumentID)
	}
	r.mu.Unlock()

	if !ok || delivery == nil {
		return nil
	}
	return delivery.Ack(false)
}

func (r *RabbitMQDocumentQueue) Nack(ctx context.Context, task queue.DocumentTask, requeue bool) error {
	r.mu.Lock()
	delivery, ok := r.deliveryMap[task.DocumentID]
	if ok {
		delete(r.deliveryMap, task.DocumentID)
	}
	r.mu.Unlock()

	if !ok || delivery == nil {
		return nil
	}
	return delivery.Nack(false, requeue)
}

func (r *RabbitMQDocumentQueue) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			logger.Error("failed to close RabbitMQ channel", err, nil)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			logger.Error("failed to close RabbitMQ connection", err, nil)
		}
	}
	return nil
}
