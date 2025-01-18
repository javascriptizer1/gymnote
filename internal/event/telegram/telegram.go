package telegram

import (
	"context"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"gymnote/internal/config"
	"gymnote/internal/entity"
)

type Fetcher interface {
	Start(ctx context.Context)
}

type fetcher struct {
	cfg        *config.TelegramConfig
	bot        *tgbotapi.BotAPI
	eventsChan chan<- entity.Event
	offset     int
}

func New(eventsChan chan<- entity.Event, cfg *config.TelegramConfig) *fetcher {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = cfg.Debug

	return &fetcher{
		cfg:        cfg,
		bot:        bot,
		eventsChan: eventsChan,
		offset:     0,
	}
}

func (f *fetcher) Start(_ context.Context) {
	u := tgbotapi.NewUpdate(f.offset)
	u.Timeout = f.cfg.Timeout

	updates := f.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil && update.Message.Text != "" && update.Message.From.ID != f.bot.Self.ID && !update.Message.IsCommand() {
			f.eventsChan <- entity.Event{
				UserID: strconv.Itoa(int(update.Message.From.ID)),
				Text:   update.Message.Text,
			}
			f.offset = update.UpdateID + 1
		}
	}
}
