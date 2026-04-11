package vk

import (
	"context"
	"encoding/json"
	"log"
	"runtime/debug"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func (a *API) Register() {
	group, err := a.vk.GroupsGetByID(api.Params{})
	if err != nil {
		log.Printf("VK: Failed to connect to group: %v", err)

		return
	}

	if len(group) > 0 {
		log.Printf("VK Bot connected to group: %s (ID: %d)", group[0].Name, group[0].ID)
	}

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

	if handler, exists := a.commandHandlers[text]; exists {
		handler(obj)
		return
	}

	if strings.HasPrefix(text, "/") {
		a.handleCommand(obj)
		return
	}

	a.handleState(obj)
}

func (a *API) handleCommand(obj *events.MessageNewObject) {
	text := strings.TrimSpace(obj.Message.Text)
	command := strings.TrimPrefix(text, "/")

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
	userID := a.getUserID(obj)
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
