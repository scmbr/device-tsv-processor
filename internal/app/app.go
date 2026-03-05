package app

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/device-tsv-processor/internal/app/config"
	delivery_http "github.com/scmbr/device-tsv-processor/internal/delivery/http"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/postgres/repository"
	"github.com/scmbr/device-tsv-processor/internal/infrastructure/rabbitmq"

	"github.com/scmbr/device-tsv-processor/internal/usecase"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type App struct {
	Config *config.Config

	DB           *sqlx.DB
	RabbitClient *rabbitmq.Client
	Repos        *repository.Repositories

	UseCases *usecase.UseCases

	Server *delivery_http.Server
}

func New(configsDir string) (*App, error) {
	cfg, err := config.Init(configsDir)
	if err != nil {
		return nil, err
	}

	db, rabbitClient, queues, repos, err := initInfrastructure(cfg)
	if err != nil {
		return nil, err
	}

	useCases := initUseCases(cfg, repos, queues)

	handler := delivery_http.NewHandler(useCases)

	server := delivery_http.NewServer(&delivery_http.ServerConfig{
		Host:               cfg.HTTP.Host,
		Port:               cfg.HTTP.Port,
		ReadTimeout:        cfg.HTTP.ReadTimeout,
		WriteTimeout:       cfg.HTTP.WriteTimeout,
		MaxHeaderMegabytes: cfg.HTTP.MaxHeaderMegabytes,
	}, handler.Init())

	return &App{
		Config:       cfg,
		DB:           db,
		RabbitClient: rabbitClient,
		Repos:        repos,
		UseCases:     useCases,
		Server:       server,
	}, nil
}
func (a *App) Run(ctx context.Context) error {

	go func() {
		if err := a.Server.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", err, nil)
		}
	}()
	go a.runWorkers(ctx)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.Shutdown(shutdownCtx)
}
func (a *App) Shutdown(ctx context.Context) error {

	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.RabbitClient.Close(); err != nil {
		return err
	}

	return a.DB.Close()
}
func (a *App) runWorkers(ctx context.Context) error {
	return nil
}
