package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	VHost     string
	QueueName string
}

type Client struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
}

func NewRabbitMQClient(cfg Config) (*Client, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		cfg.QueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		conn.Close()
		ch.Close()
		return nil, err
	}

	return &Client{
		conn:      conn,
		channel:   ch,
		QueueName: cfg.QueueName,
	}, nil
}

func (c *Client) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}
