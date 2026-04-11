package vk

import (
	"fmt"
	"strings"

	"github.com/SevereCloud/vksdk/v2/object"

	"gymnote/internal/entity"
)

func (a *API) createMuscleGroupKeyboard() *object.MessagesKeyboard {
	keyboard := &object.MessagesKeyboard{
		Inline: true,
		Buttons: [][]object.MessagesKeyboardButton{},
	}

	for _, group := range muscleGroupsWithSmiles {
		plainGroup := strings.TrimLeft(group, muscleGroupSmilePrefix)
		payload := fmt.Sprintf(`{"button":"%s"}`, musclePrefix+plainGroup)
		button := object.MessagesKeyboardButton{
			Action: object.MessagesKeyboardButtonAction{
				Type:    "callback",
				Label:   group,
				Payload: payload,
			},
		}
		keyboard.Buttons = append(keyboard.Buttons, []object.MessagesKeyboardButton{button})
	}

	return keyboard
}

func (a *API) createExerciseKeyboard(exercises []entity.Exercise, muscleGroup string, page, totalPages int, callbackPrefix string) *object.MessagesKeyboard {
	keyboard := &object.MessagesKeyboard{
		Inline: true,
		Buttons: [][]object.MessagesKeyboardButton{},
	}

	for _, exercise := range exercises {
		payload := fmt.Sprintf(`{"button":"%s"}`, callbackPrefix+exercise.ID().String())
		
		// VK ограничивает длину текста кнопки до 40 символов
		label := exercise.Name()
		if len(label) > 40 {
			label = label[:37] + "..."
		}
		
		button := object.MessagesKeyboardButton{
			Action: object.MessagesKeyboardButtonAction{
				Type:    "callback",
				Label:   label,
				Payload: payload,
			},
		}
		keyboard.Buttons = append(keyboard.Buttons, []object.MessagesKeyboardButton{button})
	}

	// Пагинация
	var paginationButtons []object.MessagesKeyboardButton
	if page > 0 {
		payload := fmt.Sprintf(`{"button":"%s%s:%s:%d"}`, musclePrefix, muscleGroup, prevDirection, page-1)
		prevButton := object.MessagesKeyboardButton{
			Action: object.MessagesKeyboardButtonAction{
				Type:    "callback",
				Label:   paginationPrevText,
				Payload: payload,
			},
		}
		paginationButtons = append(paginationButtons, prevButton)
	}

	if page < totalPages-1 {
		payload := fmt.Sprintf(`{"button":"%s%s:%s:%d"}`, musclePrefix, muscleGroup, nextDirection, page+1)
		nextButton := object.MessagesKeyboardButton{
			Action: object.MessagesKeyboardButtonAction{
				Type:    "callback",
				Label:   paginationNextText,
				Payload: payload,
			},
		}
		paginationButtons = append(paginationButtons, nextButton)
	}

	if len(paginationButtons) > 0 {
		keyboard.Buttons = append(keyboard.Buttons, paginationButtons)
	}

	return keyboard
}

func (a *API) createBackToExercisesKeyboard(muscleGroup, exerciseID string) *object.MessagesKeyboard {
	payload := fmt.Sprintf(`{"button":"%s%s:%s:0:%s"}`, musclePrefix, muscleGroup, nextDirection, exerciseID)
	keyboard := &object.MessagesKeyboard{
		Inline: true,
		Buttons: [][]object.MessagesKeyboardButton{
			{
				{
					Action: object.MessagesKeyboardButtonAction{
						Type:    "callback",
						Label:   backToExercisesText,
						Payload: payload,
					},
				},
			},
		},
	}

	return keyboard
}

func (a *API) createSetActionsKeyboard() *object.MessagesKeyboard {
	keyboard := &object.MessagesKeyboard{
		Inline: true,
		Buttons: [][]object.MessagesKeyboardButton{
			{
				{
					Action: object.MessagesKeyboardButtonAction{
						Type:    "callback",
						Label:   startNewExerciseText,
						Payload: fmt.Sprintf(`{"button":"%s"}`, startNewExercisePrefix),
					},
				},
			},
			{
				{
					Action: object.MessagesKeyboardButtonAction{
						Type:    "callback",
						Label:   finishTrainingText,
						Payload: fmt.Sprintf(`{"button":"%s"}`, confirmationFinishTrainingPrefix),
					},
				},
			},
		},
	}

	return keyboard
}

func (a *API) createFinishConfirmationKeyboard() *object.MessagesKeyboard {
	keyboard := &object.MessagesKeyboard{
		Inline: true,
		Buttons: [][]object.MessagesKeyboardButton{
			{
				{
					Action: object.MessagesKeyboardButtonAction{
						Type:    "callback",
						Label:   answerYes,
						Payload: fmt.Sprintf(`{"button":"%s"}`, acceptFinishTrainingPrefix),
					},
				},
			},
			{
				{
					Action: object.MessagesKeyboardButtonAction{
						Type:    "callback",
						Label:   answerNo,
						Payload: fmt.Sprintf(`{"button":"%s"}`, rejectFinishTrainingPrefix),
					},
				},
			},
		},
	}

	return keyboard
}
