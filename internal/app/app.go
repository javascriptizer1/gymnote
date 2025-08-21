package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gymnote/internal/chart"
	"gymnote/internal/config"
	"gymnote/internal/formatter"
	"gymnote/internal/handler/tg"
	"gymnote/internal/parser"
	"gymnote/internal/repository"
	mongodb "gymnote/internal/repository/mongo"
	"gymnote/internal/repository/redis"
	"gymnote/internal/service"
)

type app struct {
	ctx       context.Context
	cancelCtx context.CancelFunc

	cfg *config.Config

	api   tg.API
	db    repository.DB
	cache repository.Cache

	parser    service.Parser
	formatter tg.Formatter
	chart     tg.ChartService
	service   tg.TrainingService
}

func New() (*app, error) {
	ctx, cancel := context.WithCancel(context.Background())

	a := &app{
		ctx:       ctx,
		cancelCtx: cancel,
	}

	err := a.initDeps()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *app) initDeps() error {
	fns := []func() error{
		a.initConfig,
		a.initDB,
		a.initServices,
	}

	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) initConfig() error {
	cfg := config.MustLoad()
	a.cfg = cfg

	return nil
}

func (a *app) initDB() error {
	db, err := mongodb.New(a.ctx, &a.cfg.DB)
	if err != nil {
		return fmt.Errorf("init db error: %w", err)
	}

	a.db = db

	cache, err := redis.New(a.ctx, &a.cfg.Redis)
	if err != nil {
		return fmt.Errorf("init cache error: %w", err)
	}

	a.cache = cache

	return nil
}

func (a *app) initServices() error {
	a.parser = parser.New()
	a.formatter = formatter.New()
	a.chart = chart.New()
	a.service = service.New(a.db, a.cache, a.parser)

	a.api = *tg.NewAPI(a.ctx, &a.cfg.Telegram, a.formatter, a.chart, a.service)

	return nil
}

func (a *app) Run() error {
	go a.api.Register()

	log.Println("Server is running...")

	waitSignalAndShutdown(a.cancelCtx)

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.GracefulTimeout)
	defer cancel()

	return a.shutdown(ctx)
}

func (a *app) shutdown(ctx context.Context) error {
	if err := a.db.Close(ctx); err != nil {
		log.Printf("db close err: %v\n", err)
	}

	if err := a.cache.Close(ctx); err != nil {
		log.Printf("cache close err: %v\n", err)
	}

	log.Println("Shutdown complete")

	return nil
}

func waitSignalAndShutdown(cancelApp context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP, os.Interrupt,
	)

	sig := <-quit

	log.Printf("Caught signal %s. Shutting down...\n", sig)

	cancelApp()
}
