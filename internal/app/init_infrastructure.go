package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/app/config"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/repository"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq/queue"
)

func initInfrastructure(cfg *config.Config) (*sqlx.DB, *rabbitmq.Client, *queue.Queues, *repository.Repositories, error) {
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.Name,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	rabbitClient, err := rabbitmq.NewRabbitMQClient(rabbitmq.Config{
		Host:     cfg.Rabbit.Host,
		Port:     cfg.Rabbit.Port,
		Username: cfg.Rabbit.Username,
		Password: cfg.Rabbit.Password,
		VHost:    cfg.Rabbit.VHost,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	queues, err := queue.NewQueues(queue.QueuesConfig{
		Client:            rabbitClient,
		FileQueueName:     cfg.Rabbit.FileQueueName,
		DocumentQueueNAme: cfg.Rabbit.DocumentQueueName,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	repos := repository.NewRepositories(db)
	return db, rabbitClient, queues, repos, nil
}
