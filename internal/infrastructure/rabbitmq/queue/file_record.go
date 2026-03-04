package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq"
)

type RabbitFileQueue struct {
	channel   *amqp.Channel
	queueName string
}

func NewRabbitFileQueue(client *rabbitmq.Client, queueName string) (*RabbitFileQueue, error) {
	_, err := client.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &RabbitFileQueue{
		channel:   client.Channel,
		queueName: queueName,
	}, nil
}

func (q *RabbitFileQueue) Enqueue(fileID int64) error {
	body := fmt.Sprintf("%d", fileID)
	return q.channel.Publish(
		"",
		q.queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
}

func (q *RabbitFileQueue) Dequeue() (int64, error) {
	msgs, err := q.channel.Consume(
		q.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return 0, err
	}

	msg := <-msgs
	var fileID int64
	_, err = fmt.Sscanf(string(msg.Body), "%d", &fileID)
	if err != nil {
		return 0, err
	}

	return fileID, nil
}
