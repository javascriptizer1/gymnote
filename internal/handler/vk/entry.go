package vk

import (
	"context"
	"encoding/json"
	"log"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func (a *API) Register() {
	// Проверяем подключение к группе
	group, err := a.vk.GroupsGetByID(api.Params{})
	if err != nil {
		log.Printf("VK: Failed to connect to group: %v", err)
		return
	}

	if len(group) > 0 {
		log.Printf("VK Bot connected to group: %s (ID: %d)", group[0].Name, group[0].ID)
	}

	// Инициализируем Long Poll
	lp, err := longpoll.NewLongPoll(a.vk, a.cfg.GroupID)
	if err != nil {
		log.Printf("VK: Failed to initialize Long Poll: %v", err)
		return
	}

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovering from panic: %v\nStack trace: %s", r, debug.Stack())
				}
			}()

			a.handleMessage(&obj)
		}()
	})

	lp.MessageEvent(func(_ context.Context, obj events.MessageEventObject) {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovering from panic: %v\nStack trace: %s", r, debug.Stack())
				}
			}()

			a.handleCallbackQuery(&obj)
		}()
	})

	log.Println("VK LongPoll started")
	if err := lp.Run(); err != nil {
		log.Printf("VK LongPoll error: %v", err)
	}
}

func (a *API) handleMessage(obj *events.MessageNewObject) {
	text := strings.TrimSpace(obj.Message.Text)
	if text == "" {
		return
	}

	// Обработка кнопок главного меню
	switch text {
	case "🏋️ Начать тренировку":
		a.StartTrainingHandler(obj)
		return
	case "📝 Загрузить тренировку":
		a.StartUploadTrainingHandler(obj)
		return
	case "📊 История тренировок":
		a.StartGetTrainingsHandler(obj)
		return
	case "📈 Прогрессия весов":
		a.StartExerciseProgressionChartHandler(obj)
		return
	case "📋 История упражнения":
		a.StartExerciseHistoryHandler(obj)
		return
	case "🔢 Расчет 1ПМ":
		a.StartOneRMHandler(obj)
		return
	case "🗑 Сбросить тренировку":
		a.ClearTrainingHandler(obj)
		return
	case "➕ Создать упражнение":
		a.StartCreateExerciseHandler(obj)
		return
	case "❓ Помощь":
		a.HelpHandler(obj)
		return
	}

	// Проверяем, является ли сообщение командой
	if strings.HasPrefix(text, "/") {
		a.handleCommand(obj)
		return
	}

	// Обрабатываем состояние пользователя
	a.handleState(obj)
}

func (a *API) handleCommand(obj *events.MessageNewObject) {
	text := strings.TrimSpace(obj.Message.Text)
	command := strings.TrimPrefix(text, "/")
	
	// Убираем параметры команды
	if idx := strings.Index(command, " "); idx != -1 {
		command = command[:idx]
	}

	if handler, exists := a.commandHandlers[command]; exists {
		handler(obj)
	} else {
		a.UnknownCommandHandler(obj)
	}
}

func (a *API) handleState(obj *events.MessageNewObject) {
	userID := strconv.Itoa(obj.Message.FromID)
	state := a.getUserState(userID)

	if handler, exists := a.stateHandlers[state]; exists {
		handler(obj)
	} else {
		a.UnknownCommandHandler(obj)
	}
}

func (a *API) handleCallbackQuery(obj *events.MessageEventObject) {
	payload := parseCallbackPayload(obj.Payload)
	if payload == "" {
		return
	}

	for prefix, handler := range a.callbackHandlers {
		if strings.HasPrefix(payload, prefix) {
			handler(obj)
			return
		}
	}
}


// parseCallbackPayload извлекает данные кнопки из JSON payload
func parseCallbackPayload(payload []byte) string {
	var payloadData struct {
		Button string `json:"button"`
	}
	
	if err := json.Unmarshal(payload, &payloadData); err != nil {
		log.Printf("VK: Failed to parse callback payload: %v", err)
		return ""
	}
	
	return payloadData.Button
}
