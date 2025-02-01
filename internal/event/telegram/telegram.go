package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"

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
		if update.CallbackQuery != nil {
			handleCallback(f.bot, update.CallbackQuery)
		}
		if update.Message != nil && update.Message.IsCommand() {
			StartTrainingHandler(f.bot, update)
		}
		if update.Message != nil && update.Message.Text != "" && update.Message.From.ID != f.bot.Self.ID && !update.Message.IsCommand() {
			f.eventsChan <- entity.Event{
				UserID: strconv.Itoa(int(update.Message.From.ID)),
				Text:   update.Message.Text,
			}
			f.offset = update.UpdateID + 1
		}
	}
}

func StartTrainingHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Спина", "muscle:spina"),
			tgbotapi.NewInlineKeyboardButtonData("Грудь", "muscle:grud"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ноги", "muscle:nogi"),
			tgbotapi.NewInlineKeyboardButtonData("Бицепс", "muscle:biceps"),
		),
	)

	// Сохранить начальное состояние в Redis
	// initialState := map[string]interface{}{
	// 	"session_id":       uuid.New().String(),
	// 	"current_muscle_group": "",
	// 	"current_exercise": "",
	// 	"sets":             []map[string]interface{}{},
	// }
	// SaveSessionState(update.Message.From.ID, initialState)

	// Отправить сообщение с кнопками
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите мышечную группу:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func MuscleGroupHandler(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Получить мышечную группу из callback data
	// muscleGroup := strings.Split(callbackQuery.Data, ":")[1] // Например, "spina"

	// Обновить состояние
	// state, _ := GetSessionState(callbackQuery.From.ID)
	// state["current_muscle_group"] = muscleGroup
	// SaveSessionState(callbackQuery.From.ID, state)

	// Отправить сообщение с кнопками упражнений
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тяга штанги", "exercise:tyaga"),
			tgbotapi.NewInlineKeyboardButtonData("Подтягивания", "exercise:podtyag"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тяга верхнего блока", "exercise:blocco"),
		),
	)

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Выберите упражнение:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func handleCallback(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Проверяем данные из callback
	callbackData := callbackQuery.Data

	// Фильтруем по префиксу, например: "muscle:*"
	if strings.HasPrefix(callbackData, "muscle:") {
		MuscleGroupHandler(bot, callbackQuery) // Вызываем хендлер для выбора мышечной группы
	} else if strings.HasPrefix(callbackData, "exercise:") {
		// ExerciseHandler(bot, callbackQuery) // Вызываем хендлер для выбора упражнения
	} else if callbackData == "finish_training" {
		// FinishTrainingHandler(bot, callbackQuery) // Завершение тренировки
	} else if callbackData == "choose_new" {
		// Пример нового действия
		// StartTrainingHandler(bot, callbackQuery.Message) // Возвращаемся к выбору
	}

	// После обработки callback обязательно отправляем ответ, чтобы убрать "загрузка" в интерфейсе Telegram
	// bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackQuery.ID, ""))
}
