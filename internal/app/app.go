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
	"golang.org/x/sync/errgroup"

	"github.com/scmbr/device-tsv-processor/internal/usecase"
	"github.com/scmbr/device-tsv-processor/pkg/logger"
)

type App struct {
	Config       *config.Config
	DB           *sqlx.DB
	RabbitClient *rabbitmq.Client
	Repos        *repository.Repositories

	UseCases *usecase.UseCases
	Workers  *Workers
	Server   *delivery_http.Server
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
	workers := initWorkers(cfg, useCases, queues)
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
		Workers:      workers,
		Server:       server,
	}, nil
}
func (a *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := a.Server.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", err, nil)
			return err
		}
		return nil
	})

	a.startWorkers(g, ctx)

	err := g.Wait()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if shutErr := a.Shutdown(shutdownCtx); shutErr != nil {
		logger.Error("shutdown failed", shutErr, nil)
	}

	return err
}

func (a *App) startWorkers(g *errgroup.Group, ctx context.Context) {

	g.Go(func() error {
		return a.Workers.Generator.StartPool(ctx, a.Config.Workers.GeneratorWorkersCount)
	})

	g.Go(func() error {
		return a.Workers.Process.StartPool(ctx, a.Config.Workers.ProcessWorkersCount)
	})

	g.Go(func() error {
		return a.Workers.Scan.Start(ctx)
	})

	g.Go(func() error {
		return a.Workers.Queue.Start(ctx)
	})
}

func (a *App) Shutdown(ctx context.Context) error {
	var firstErr error

	if err := a.Server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", err, nil)
		if firstErr == nil {
			firstErr = err
		}
	}
	logger.Info("server stopped", nil)
	if err := a.RabbitClient.Close(); err != nil {
		logger.Error("rabbit client close failed", err, nil)
		if firstErr == nil {
			firstErr = err
		}
	}
	logger.Info("rabbit stopped", nil)
	if err := a.DB.Close(); err != nil {
		logger.Error("db close failed", err, nil)
		if firstErr == nil {
			firstErr = err
		}
	}
	logger.Info("db stopped", nil)
	return firstErr
}
