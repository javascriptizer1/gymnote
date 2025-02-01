package tg

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *API) Register(updates tgbotapi.UpdatesChannel) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovering from panic: %v\nStack trace: %s", r, debug.Stack())
		}
	}()

	for update := range updates {
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
	commandHandlers := map[string]func(*tgbotapi.Message){
		"/start_training": a.StartTrainingHandler,
	}

	if handler, exists := commandHandlers[message.Text]; exists {
		handler(message)
	}
}

func (a *API) handleState(message *tgbotapi.Message) {
	userID := strconv.FormatInt(message.From.ID, 10)
	session, err := a.trainingService.GetCurrentSession(a.ctx, userID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		a.bot.Send(msg)
		return
	}

	switch session.State() {
	case "awaiting_set_input":
		a.SetHandler(message)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда.")
		a.bot.Send(msg)
	}
}

func (a *API) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	callbackHandlers := map[string]func(*tgbotapi.CallbackQuery){
		"muscle:":         a.MuscleGroupHandler,
		"exercise:":       a.ExerciseHandler,
		"finish_training": a.FinishTrainingHandler,
	}

	for prefix, handler := range callbackHandlers {
		if strings.HasPrefix(callback.Data, prefix) {
			handler(callback)
			return
		}
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Неизвестная команда.")
	a.bot.Send(msg)
}
