package tg

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"gymnote/internal/chart"
	"gymnote/internal/entity"
	"gymnote/internal/errs"
)

var (
	pageSize  = 5
	parseMode = tgbotapi.ModeMarkdown
)

func (a *API) StartHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	if a.cfg.GreetingStickerID != "" {
		sticker := tgbotapi.NewSticker(chatID, tgbotapi.FileID(a.cfg.GreetingStickerID))
		_, _ = a.bot.Send(sticker)
	}

	text := startText
	if a.cfg.AuthorName != "" {
		text += fmt.Sprintf(donateAuthorText, a.cfg.AuthorName)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = a.bot.Send(msg)
}

func (a *API) HelpHandler(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	_, _ = a.bot.Send(msg)
}

func (a *API) StartExerciseProgressionChartHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, group := range muscleGroupsWithSmiles {
		plainGroup := strings.TrimLeft(group, "💪🏋️🦵 ")
		button := tgbotapi.NewInlineKeyboardButtonData(group, musclePrefix+plainGroup)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(button))
	}

	msg := tgbotapi.NewMessage(chatID, startProgressionMuscleGroupSelectText)
	msg.ParseMode = parseMode
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	a.setUserState(userID, entity.StateAwaitingExerciseProgression)

	_, _ = a.bot.Send(msg)
}

func (a *API) ExerciseProgressionChartHandler(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	userID := strconv.FormatInt(callback.From.ID, 10)
	messageID := callback.Message.MessageID
	exerciseIDStr := strings.TrimPrefix(callback.Data, startGetExerciseProgressionPrefix)

	defer a.clearUserState(userID)

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, errInvalidExerciseID)
		_, _ = a.bot.Send(msg)
		return
	}

	loadingMsg := tgbotapi.NewEditMessageText(chatID, messageID, loadingProgressionText)
	_, _ = a.bot.Send(loadingMsg)

	data, err := a.trainingService.GetExerciseProgression(a.ctx, userID, exerciseID)
	if err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, errGetTrainings))
		return
	}
	if len(data) == 0 {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, notFoundTrainingsText))
		return
	}

	var xValues []string
	var yValues []float32

	for _, v := range data {
		xValues = append(xValues, v.SessionDate.Format("2006-01-02"))
		yValues = append(yValues, v.Weight)
	}

	exerciseName := data[0].ExerciseName
	cfg := chart.LinearChartConfig{
		Title:    exerciseName,
		XName:    "Дата",
		YName:    "Вес (кг)",
		YValues:  yValues,
		XValues:  xValues,
		FileName: fmt.Sprintf("%s/%s-%s.png", a.cfg.GraphicsPath, userID, time.Now().Format("2006-01-02")),
	}

	if err = a.chartService.GenerateLinearChart(cfg); err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, errProgression))
		return
	}

	chartImage := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(cfg.FileName))
	_, _ = a.bot.Send(chartImage)
}

func (a *API) StartGetTrainingsHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	a.setUserState(userID, entity.StateAwaitingGetTrainingsInput)
	_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, startGetTrainingsText))
}

func (a *API) GetTrainingsHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	defer a.clearUserState(userID)

	var fromDate, toDate *time.Time
	args := strings.SplitN(message.Text, " ", 2)

	if len(args) >= 1 {
		from, err := time.Parse("2006-01-02", args[0])
		if err == nil {
			fromDate = &from
		}
	}
	if len(args) == 2 {
		to, err := time.Parse("2006-01-02", args[1])
		if err == nil {
			toDate = &to
		}
	}

	trainings, err := a.trainingService.GetTrainingSessions(a.ctx, userID, fromDate, toDate)
	if err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, errGetTrainings))
		return
	}
	if len(trainings) == 0 {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, notFoundTrainingsText))
		return
	}

	text := a.formatter.FormatTrainingLogs(trainings)
	_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, text))
}

func (a *API) UnknownCommandHandler(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, unknownCommandText)
	_, _ = a.bot.Send(msg)
}

func (a *API) StartCreateExerciseHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	a.setUserState(userID, entity.StateAwaitingExerciseInput)
	_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, startCreateExerciseText))
}

func (a *API) CreateExerciseHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	defer a.clearUserState(userID)

	args := strings.SplitN(message.Text, " ", 3)
	if len(args) < 3 {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf(errGeneral, startCreateExerciseText)))
		return
	}

	name, muscleGroup, equipment := args[0], args[1], args[2]
	if name == "" {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf(errGeneral, emptyExerciseNameText)))
		return
	}
	if !slices.Contains(availableMuscleGroups, muscleGroup) {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf(errGeneral, fmt.Sprintf(unknownMuscleGroupText, availableMuscleGroups))))
		return
	}

	err := a.trainingService.CreateExercise(a.ctx, name, muscleGroup, equipment)
	if err != nil {
		if errors.Is(err, errs.ErrExerciseAlreadyExists) {
			_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf(errGeneral, fmt.Sprintf(exerciseWithNameAlreadyExistsText, name))))
		} else {
			_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, errCreateExercise))
		}
		return
	}

	_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf(exerciseCreatedText, name, muscleGroup)))
}

func (a *API) StartUploadTrainingHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	a.setUserState(userID, entity.StateAwaitingTrainingInput)
	_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, startUploadTrainingText))
}

func (a *API) UploadTrainingHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	defer a.clearUserState(userID)

	session, err := a.trainingService.ParseTraining(a.ctx, entity.Event{UserID: userID, Text: message.Text})
	if err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, errUploadTraining))
		return
	}

	text := fmt.Sprintf(finishText, session.ExerciseCount(), session.SetCount(), session.TotalVolume())
	_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, text))
}

func (a *API) ClearTrainingHandler(message *tgbotapi.Message) {
	userID := strconv.FormatInt(message.From.ID, 10)

	defer a.clearUserState(userID)

	if err := a.trainingService.ClearSession(a.ctx, userID); err != nil {
		text := errClearTraining
		if errors.Is(err, errs.ErrSessionNotFound) {
			text = errNoTraining
		}
		msg := tgbotapi.NewMessage(message.From.ID, text)
		_, _ = a.bot.Send(msg)
		return
	}

	_, _ = a.bot.Send(tgbotapi.NewMessage(message.From.ID, clearTrainingDoneText))
}

func (a *API) StartTrainingHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	_, err := a.trainingService.StartTraining(a.ctx, userID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(errStartTraining, err))
		_, _ = a.bot.Send(msg)
		return
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, group := range muscleGroupsWithSmiles {
		plainGroup := strings.TrimLeft(group, "💪🏋️🦵 ")
		button := tgbotapi.NewInlineKeyboardButtonData(group, musclePrefix+plainGroup)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(button))
	}

	msg := tgbotapi.NewMessage(chatID, startTrainingText)
	msg.ParseMode = parseMode
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	_, _ = a.bot.Send(msg)
}

func (a *API) MuscleGroupHandler(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	userID := strconv.FormatInt(callback.From.ID, 10)
	state := a.getUserState(userID)

	muscleGroup, page, _, err := parseMuscleGroupCallbackData(callback.Data)
	if err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, errGeneral))
		return
	}

	exercises, err := a.trainingService.GetExercisesByMuscleGroup(a.ctx, muscleGroup)
	if err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, errExerciseLoad))
		return
	}

	pagedExercises, totalPages, err := paginate(exercises, page, pageSize)
	if err != nil {
		_, _ = a.bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		return
	}

	callbackDataPrefix := exercisePrefix
	if state == entity.StateAwaitingExerciseProgression {
		callbackDataPrefix = startGetExerciseProgressionPrefix
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, exercise := range pagedExercises {
		button := tgbotapi.NewInlineKeyboardButtonData(exercise.Name(), callbackDataPrefix+exercise.ID().String())
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(button))
	}

	var paginationButtons []tgbotapi.InlineKeyboardButton
	if page > 0 {
		prevButton := tgbotapi.NewInlineKeyboardButtonData(paginationPrevText, fmt.Sprintf("%s%s:%s:%d", musclePrefix, muscleGroup, prevDirection, page-1))
		paginationButtons = append(paginationButtons, prevButton)
	}

	if page < totalPages-1 {
		nextButton := tgbotapi.NewInlineKeyboardButtonData(paginationNextText, fmt.Sprintf("%s%s:%s:%d", musclePrefix, muscleGroup, nextDirection, page+1))
		paginationButtons = append(paginationButtons, nextButton)
	}

	if len(paginationButtons) > 0 {
		buttons = append(buttons, paginationButtons)
	}

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf(muscleGroupDoneText, muscleGroup))
	editMsg.ParseMode = parseMode

	editMarkup := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, tgbotapi.NewInlineKeyboardMarkup(buttons...))

	_, _ = a.bot.Send(editMsg)
	_, _ = a.bot.Send(editMarkup)
}

func (a *API) ExerciseHandler(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	userID := strconv.FormatInt(callback.From.ID, 10)
	exerciseIDStr := strings.TrimPrefix(callback.Data, exercisePrefix)

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, errInvalidExerciseID)
		_, _ = a.bot.Send(msg)
		return
	}

	err = a.trainingService.AddTrainingExercise(a.ctx, userID, exerciseID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(errAddExercise, err))
		_, _ = a.bot.Send(msg)
		return
	}

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, exerciseText)
	editMsg.ParseMode = parseMode

	a.setUserState(userID, entity.StateAwaitingSetInput)
	_, _ = a.bot.Send(editMsg)
}

func (a *API) SetHandler(message *tgbotapi.Message) {
	userID := strconv.FormatInt(message.From.ID, 10)
	input := message.Text

	parts := strings.SplitN(input, "\n", 2)
	setData := strings.Split(parts[0], ",")
	if len(setData) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, errInvalidFormat)
		_, _ = a.bot.Send(msg)
		return
	}

	weight, errWeight := strconv.ParseFloat(setData[0], 64)
	reps, errReps := strconv.Atoi(setData[1])
	if errWeight != nil || errReps != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, errParseData)
		_, _ = a.bot.Send(msg)
		return
	}
	var notes string
	if len(parts) > 1 {
		notes = strings.TrimSpace(parts[1])
	}

	err := a.trainingService.AddOrUpdateSet(a.ctx, userID, float32(weight), uint8(reps), notes)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(errGeneral, err))
		_, _ = a.bot.Send(msg)
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(startNewExerciseText, startNewExercisePrefix),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(finishTrainingText, finishTrainingPrefix),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, setText)
	msg.ReplyMarkup = keyboard

	_, _ = a.bot.Send(msg)
}

func (a *API) StartNewExerciseHandler(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	userID := strconv.FormatInt(callback.From.ID, 10)

	defer a.clearUserState(userID)

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, group := range muscleGroupsWithSmiles {
		plainGroup := strings.TrimLeft(group, "💪🏋️🦵 ")
		button := tgbotapi.NewInlineKeyboardButtonData(group, musclePrefix+plainGroup)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(button))
	}

	msg := tgbotapi.NewMessage(chatID, muscleGroupSelectText)
	msg.ParseMode = parseMode
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	_, _ = a.bot.Send(msg)
}

func (a *API) FinishTrainingHandler(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	userID := strconv.FormatInt(callback.From.ID, 10)

	defer a.clearUserState(userID)

	session, err := a.trainingService.EndSession(a.ctx, userID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(errGeneral, err))
		_, _ = a.bot.Send(msg)
		return
	}

	text := fmt.Sprintf(finishText, session.ExerciseCount(), session.SetCount(), session.TotalVolume())
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = parseMode

	_, _ = a.bot.Send(editMsg)
}

func parseMuscleGroupCallbackData(data string) (muscleGroup string, page int, direction string, err error) {
	parts := strings.Split(data, ":")

	if len(parts) < 2 {
		return "", 0, "", fmt.Errorf("invalid callback format")
	}

	muscleGroup = parts[1]
	page = 0
	direction = ""

	if len(parts) == 4 {
		direction = parts[2]
		page, err = strconv.Atoi(parts[3])
		if err != nil {
			return "", 0, "", fmt.Errorf("invalid page number")
		}
	}

	return muscleGroup, page, direction, nil
}

func paginate[T any](entities []T, page int, pageSize int) ([]T, int, error) {
	if page < 0 || pageSize <= 0 {
		return nil, 0, fmt.Errorf("invalid page number or page size")
	}

	totalPages := (len(entities) + pageSize - 1) / pageSize

	if page >= totalPages {
		return nil, 0, fmt.Errorf("page number exceeds total pages")
	}

	start := page * pageSize
	end := start + pageSize

	if end > len(entities) {
		end = len(entities)
	}

	pagedExercises := entities[start:end]
	return pagedExercises, totalPages, nil
}
