package tg

import (
	"log"
	"runtime/debug"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *API) Register() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovering from panic: %v\nStack trace: %s", r, debug.Stack())
		}
	}()

	for update := range a.bot.GetUpdatesChan(tgbotapi.UpdateConfig{}) {
		switch {
		case update.Message != nil && update.Message.IsCommand():
			a.handleCommand(update.Message)
		case update.Message != nil:
			a.handleState(update.Message)
		case update.CallbackQuery != nil:
			a.handleCallbackQuery(update.CallbackQuery)
		default:
			log.Printf("Unknown update type: %+v\n", update)
		}
	}
}

func (a *API) handleCommand(message *tgbotapi.Message) {
	if handler, exists := a.commandHandlers[message.Command()]; exists {
		handler(message)
	} else {
		a.UnknownCommandHandler(message)
	}
}

func (a *API) handleState(message *tgbotapi.Message) {
	userID := strconv.FormatInt(message.From.ID, 10)
	state := a.getUserState(userID)

	if handler, exists := a.stateHandlers[state]; exists {
		handler(message)
	} else {
		a.UnknownCommandHandler(message)
	}
}

func (a *API) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	for prefix, handler := range a.callbackHandlers {
		if strings.HasPrefix(callback.Data, prefix) {
			handler(callback)
			return
		}
	}

	a.UnknownCommandHandler(callback.Message)
}
