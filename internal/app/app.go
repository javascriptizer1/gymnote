package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gymnote/internal/config"
	"gymnote/internal/consumer"
	"gymnote/internal/entity"
	"gymnote/internal/event/telegram"
	"gymnote/internal/parser"
	"gymnote/internal/repository"
	"gymnote/internal/repository/clickhouse"
	"gymnote/internal/service"
)

type app struct {
	cfg *config.Config

	db repository.DB

	fetcher   telegram.Fetcher
	parser    service.Parser
	processor consumer.Processor
	consumer  consumer.Consumer

	eventChan chan entity.Event
	wg        *sync.WaitGroup
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func New() (*app, error) {
	ctx, cancel := context.WithCancel(context.Background())

	a := &app{
		wg:        &sync.WaitGroup{},
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
		a.initChan,
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
	db, err := clickhouse.New(a.ctx, &a.cfg.DB)
	if err != nil {
		return fmt.Errorf("init db error: %w", err)
	}

	a.db = db

	return nil
}

func (a *app) initChan() error {
	a.eventChan = make(chan entity.Event, 100)
	return nil
}

func (a *app) initServices() error {
	a.fetcher = telegram.New(a.eventChan, &a.cfg.Telegram)
	a.parser = parser.New()
	a.processor = service.New(a.db, a.parser)
	a.consumer = consumer.New(a.eventChan, a.processor)

	return nil
}

func (a *app) Run() error {
	defer a.shutdown()
	return a.runServer()
}

func (a *app) shutdown() {
	a.db.Close()
}

func (a *app) runServer() error {
	go a.fetcher.Start(a.ctx)
	go a.consumer.Start(a.ctx)

	initSignalHandler(a.wg, a.cancelCtx)

	return nil
}

func initSignalHandler(wg *sync.WaitGroup, cancel context.CancelFunc) {
	osSigCh := make(chan os.Signal, 1)
	signal.Notify(osSigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	wg.Add(1)

	go func() {
		defer wg.Done()

		signalReceived := <-osSigCh

		switch signalReceived {
		case syscall.SIGINT:
			log.Println("Received SIGINT, initiating graceful shutdown...")
		case syscall.SIGTERM:
			log.Println("Received SIGTERM, initiating graceful shutdown...")
		case syscall.SIGQUIT:
			log.Println("Received SIGQUIT, initiating graceful shutdown...")
		default:
			log.Println("Received unknown signal, initiating graceful shutdown...")
		}

		cancel()

		time.Sleep(2 * time.Second)
	}()

	wg.Wait()
}
