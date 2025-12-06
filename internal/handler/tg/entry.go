package tg

import (
	"log"
	"runtime/debug"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *API) Register() {
	for update := range a.bot.GetUpdatesChan(tgbotapi.UpdateConfig{}) {
		go func(update tgbotapi.Update) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovering from panic: %v\nStack trace: %s", r, debug.Stack())
				}
			}()

			switch {
			case update.EditedMessage != nil:
				a.handleEditedMessage(update.EditedMessage)
			case update.Message != nil && update.Message.IsCommand():
				a.handleCommand(update.Message)
			case update.Message != nil:
				a.handleState(update.Message)
			case update.CallbackQuery != nil:
				a.handleCallbackQuery(update.CallbackQuery)
			default:
				log.Printf("Unknown update type: %+v\n", update)
			}
		}(update)
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

func (a *API) handleEditedMessage(message *tgbotapi.Message) {
	if message == nil || message.From == nil {
		return
	}

	userID := strconv.FormatInt(message.From.ID, 10)
	input := message.Text
	if input == "" {
		return
	}

	parts := strings.SplitN(input, "\n", 2)
	setData := strings.Split(parts[0], ",")
	if len(setData) != 2 {
		return
	}

	weight, errWeight := strconv.ParseFloat(setData[0], 64)
	reps, errReps := strconv.Atoi(setData[1])
	if errWeight != nil || errReps != nil {
		return
	}

	var notes string
	if len(parts) > 1 {
		notes = strings.TrimSpace(parts[1])
	}

	_ = a.trainingService.UpdateSetFromMessage(a.ctx, userID, message.MessageID, float32(weight), uint8(reps), notes)
}
