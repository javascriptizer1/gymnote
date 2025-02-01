package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"gymnote/internal/config"
	"gymnote/internal/consumer"
	"gymnote/internal/entity"
	"gymnote/internal/event/telegram"
	"gymnote/internal/handler/tg"
	"gymnote/internal/parser"
	"gymnote/internal/repository"
	"gymnote/internal/repository/clickhouse"
	"gymnote/internal/repository/redis"
	"gymnote/internal/service"
)

type app struct {
	cfg *config.Config

	bot *tgbotapi.BotAPI
	api tg.API

	db    repository.DB
	cache repository.Cache

	fetcher   telegram.Fetcher
	parser    service.Parser
	processor tg.TrainingService
	consumer  consumer.Consumer

	eventChan chan entity.Event
	ctx       context.Context
	cancelCtx context.CancelFunc
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

	cache, err := redis.New(a.ctx, &a.cfg.Redis)
	if err != nil {
		return fmt.Errorf("init cache error: %w", err)
	}

	a.cache = cache

	return nil
}

func (a *app) initChan() error {
	a.eventChan = make(chan entity.Event, 100)
	return nil
}

func (a *app) initServices() error {

	bot, err := tgbotapi.NewBotAPI(a.cfg.Telegram.BotToken)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = a.cfg.Telegram.Debug

	a.bot = bot

	a.fetcher = telegram.New(a.eventChan, &a.cfg.Telegram)
	a.parser = parser.New()
	a.processor = service.New(a.db, a.cache, a.parser)
	// a.consumer = consumer.New(a.eventChan, a.processor)

	a.api = *tg.New(a.ctx, a.bot, a.processor)

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
	// go a.fetcher.Start(a.ctx)
	// go a.consumer.Start(a.ctx)

	a.api.Register(a.bot.GetUpdatesChan(tgbotapi.UpdateConfig{}))

	waitGracefulShutdown(a.cancelCtx, a.cfg.GracefulTimeout)

	return nil
}

func waitGracefulShutdown(cancel context.CancelFunc, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGHUP, os.Interrupt,
	)

	log.Printf("Caught signal %s. Shutting down...\n", <-quit)

	done := make(chan struct{})
	go func() {
		cancel()
		done <- struct{}{}
	}()

	select {
	case <-time.After(timeout):
	case <-done:
	}
}
