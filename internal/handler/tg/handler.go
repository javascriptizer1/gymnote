package tg

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

func (a *API) StartTrainingHandler(message *tgbotapi.Message) {
	userID := strconv.FormatInt(message.From.ID, 10)

	_, err := a.trainingService.StartTraining(a.ctx, userID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		a.bot.Send(msg)
		return
	}

	muscleGroups := []string{"Спина", "Грудь", "Ноги", "Руки"}
	rows := make([]tgbotapi.InlineKeyboardButton, len(muscleGroups))
	for i, group := range muscleGroups {
		rows[i] = tgbotapi.NewInlineKeyboardButtonData(group, "muscle:"+group)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите мышечную группу:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows)
	a.bot.Send(msg)
}

func (a *API) MuscleGroupHandler(callback *tgbotapi.CallbackQuery) {
	group := strings.TrimPrefix(callback.Data, "muscle:")
	exercises, err := a.trainingService.GetExercisesByMuscleGroup(a.ctx, group)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("Ошибка при загрузке упражнений: %v", err))
		a.bot.Send(msg)
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(exercises); i += 1 {
		end := i + 1
		if end > len(exercises) {
			end = len(exercises)
		}
		page := exercises[i:end]
		buttons := make([]tgbotapi.InlineKeyboardButton, len(page))
		for j, exercise := range page {
			buttons[j] = tgbotapi.NewInlineKeyboardButtonData(
				shortenExerciseName(exercise.Name()),
				fmt.Sprintf("exercise:%s", exercise.ID().String()),
			)
		}
		rows = append(rows, buttons)
	}
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("Вы выбрали группу: %s. Теперь выберите упражнение:", group))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	a.bot.Send(msg)
}

func shortenExerciseName(name string) string {
	// if len(name) > 15 { // Ограничиваем длину названия
	// 	return name[:15] + "..." // Обрезаем и добавляем многоточие
	// }
	return name
}

func (a *API) ExerciseHandler(callback *tgbotapi.CallbackQuery) {
	data := strings.TrimPrefix(callback.Data, "exercise:")
	// parts := strings.SplitN(data, ":", 2)
	// if len(parts) != 2 {
	// 	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Неверные данные упражнения.")
	// 	a.bot.Send(msg)
	// 	return
	// }

	// exerciseIDStr, exerciseName := parts[0], parts[1]
	exerciseID, err := uuid.Parse(data)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ошибка: неверный формат ID упражнения.")
		a.bot.Send(msg)
		return
	}

	userID := strconv.FormatInt(callback.From.ID, 10)
	err = a.trainingService.AddExercise(a.ctx, userID, exerciseID)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		a.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("Вы выбрали упражнение: %s. Введите вес и повторения через запятую:", data))
	a.bot.Send(msg)
}

func (a *API) SetHandler(message *tgbotapi.Message) {
	input := message.Text
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неверный формат. Введите вес и повторения через запятую (например: 50.5,12):")
		a.bot.Send(msg)
		return
	}

	userID := strconv.FormatInt(message.From.ID, 10)
	weight, errWeight := strconv.ParseFloat(parts[0], 64)
	reps, errReps := strconv.Atoi(parts[1])
	if errWeight != nil || errReps != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при разборе данных. Проверьте формат и попробуйте снова.")
		a.bot.Send(msg)
		return
	}

	err := a.trainingService.UpdateActiveSet(a.ctx, userID, float32(weight), uint8(reps))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		a.bot.Send(msg)
		return
	}

	s, err := a.trainingService.GetCurrentSession(a.ctx, userID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		a.bot.Send(msg)
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить заметку", fmt.Sprintf("note:%s", s.ActiveExercise().Exercise.ID())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить новое упражнение", "start_training"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Завершить тренировку", "finish_training"),
		),
	)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Подход сохранен! Введите данные нового подхода, либо:")
	msg.ReplyMarkup = keyboard
	a.bot.Send(msg)
}

func (a *API) FinishTrainingHandler(callback *tgbotapi.CallbackQuery) {
	userID := strconv.FormatInt(callback.From.ID, 10)

	session, err := a.trainingService.EndSession(a.ctx, userID)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		a.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf(`
	Тренировка завершена! Все данные сохранены.
	Сегодня вы:
	- Выполнили %d упражнений
	- Сделали %d подходов
	- Подняли %.2f кг в сумме
	`,
		session.ExerciseCount(),
		session.SetCount(),
		session.TotalVolume()),
	)
	a.bot.Send(msg)
}
