package vk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/google/uuid"

	"gymnote/internal/chart"
	"gymnote/internal/entity"
	"gymnote/internal/errs"
	"gymnote/internal/formatter"
	"gymnote/internal/handler/common"
	"gymnote/internal/onerm"
)

var (
	maxVKMessageLength     = 4096
	daysForSetStatistics   = common.DaysForSetStatistics
	daysForExerciseHistory = common.DaysForExerciseHistory
	pageSize               = 5
)

func (a *API) sendMessage(peerID int, message string, keyboard *object.MessagesKeyboard) {
	b := params.NewMessagesSendBuilder()
	b.PeerID(peerID)
	b.Message(message)
	b.RandomID(0)

	if keyboard != nil {
		b.Keyboard(keyboard)
	}

	_, _ = a.vk.MessagesSend(b.Params)
}

func (a *API) sendPhoto(peerID int, photoPath string) {
	uploadServer, err := a.vk.PhotosGetMessagesUploadServer(api.Params{
		"peer_id": peerID,
	})
	if err != nil {
		log.Printf("VK: Failed to get upload server: %v", err)
		return
	}

	fileData, err := os.ReadFile(photoPath)
	if err != nil {
		log.Printf("VK: Failed to read file: %v", err)
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("photo", filepath.Base(photoPath))
	if err != nil {
		log.Printf("VK: Failed to create form file: %v", err)
		return
	}

	if _, err := part.Write(fileData); err != nil {
		log.Printf("VK: Failed to write file data: %v", err)
		return
	}

	if err := writer.Close(); err != nil {
		log.Printf("VK: Failed to close writer: %v", err)
		return
	}

	req, err := http.NewRequest("POST", uploadServer.UploadURL, body)
	if err != nil {
		log.Printf("VK: Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("VK: Failed to upload: %v", err)
		return
	}
	defer resp.Body.Close()

	uploadResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("VK: Failed to read upload response: %v", err)
		return
	}

	var uploadResult struct {
		Photo  string `json:"photo"`
		Server int    `json:"server"`
		Hash   string `json:"hash"`
	}

	if err := json.Unmarshal(uploadResp, &uploadResult); err != nil {
		log.Printf("VK: Failed to parse upload response: %v, raw: %s", err, string(uploadResp))
		return
	}

	if uploadResult.Photo == "" || uploadResult.Photo == "[]" {
		log.Printf("VK: Upload failed - empty photo data. Response: %s", string(uploadResp))
		return
	}

	photos, err := a.vk.PhotosSaveMessagesPhoto(api.Params{
		"photo":  uploadResult.Photo,
		"server": uploadResult.Server,
		"hash":   uploadResult.Hash,
	})
	if err != nil {
		log.Printf("VK: Failed to save photo: %v", err)
		return
	}

	if len(photos) == 0 {
		log.Printf("VK: No photos returned")
		return
	}

	attachment := fmt.Sprintf("photo%d_%d", photos[0].OwnerID, photos[0].ID)

	b := params.NewMessagesSendBuilder()
	b.PeerID(peerID)
	b.Attachment(attachment)
	b.RandomID(0)

	_, err = a.vk.MessagesSend(b.Params)
	if err != nil {
		log.Printf("VK: Failed to send photo: %v", err)
	}
}

func (a *API) editMessage(peerID, conversationMessageID int, message string, keyboard *object.MessagesKeyboard) {
	b := params.NewMessagesEditBuilder()
	b.PeerID(peerID)
	b.ConversationMessageID(conversationMessageID)
	b.Message(message)

	if keyboard != nil {
		b.Keyboard(keyboard)
	}

	_, _ = a.vk.MessagesEdit(b.Params)
}

func (a *API) sendEventAnswer(eventID string, userID int, peerID int) {
	b := params.NewMessagesSendMessageEventAnswerBuilder()
	b.EventID(eventID)
	b.UserID(userID)
	b.PeerID(peerID)

	_, _ = a.vk.MessagesSendMessageEventAnswer(b.Params)
}

func (a *API) StartHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID

	text := startText
	if a.cfg.AuthorName != "" {
		text += fmt.Sprintf(donateAuthorText, a.cfg.AuthorName)
	}

	keyboard := a.createMainMenuKeyboard()

	a.sendMessage(peerID, text, keyboard)
}

func (a *API) HelpHandler(obj *events.MessageNewObject) {
	keyboard := a.createMainMenuKeyboard()
	a.sendMessage(obj.Message.PeerID, helpText, keyboard)
}

func (a *API) OneRMHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	defer a.clearUserState(userID)

	input := strings.TrimSpace(obj.Message.Text)
	parts := strings.SplitN(input, "\n", 2)
	setData := strings.Split(parts[0], ",")
	if len(setData) != 2 {
		a.sendMessage(peerID, errInvalidFormat, nil)

		return
	}

	weight, errWeight := strconv.ParseFloat(strings.TrimSpace(setData[0]), 64)
	reps, errReps := strconv.Atoi(strings.TrimSpace(setData[1]))
	if errWeight != nil || errReps != nil || weight <= 0 || reps <= 0 {
		a.sendMessage(peerID, errParseData, nil)

		return
	}

	summary := onerm.Calculate(weight, reps)
	if len(summary.Results) == 0 || summary.Average <= 0 {
		a.sendMessage(peerID, errGeneral, nil)

		return
	}

	formulaNames := map[onerm.Formula]string{
		onerm.FormulaEpley:    "Эпли",
		onerm.FormulaBrzycki:  "Бжицки",
		onerm.FormulaLander:   "Лэндер",
		onerm.FormulaLombardi: "Ломбарди",
		onerm.FormulaMayhew:   "Мэйхью",
		onerm.FormulaOConner:  "О'Коннор",
		onerm.FormulaWathan:   "Ватан",
	}

	var sb strings.Builder
	sb.WriteString("📈 Расчёт одноповторного максимума\n\n")
	sb.WriteString(fmt.Sprintf("Исходные данные: %s кг x %d\n\n", formatter.FormatWeightFloat(weight), reps))
	sb.WriteString("Формулы:\n\n")
	for _, r := range summary.Results {
		name := formulaNames[r.Formula]
		if name == "" {
			name = string(r.Formula)
		}
		sb.WriteString(fmt.Sprintf("• %s: %s кг\n", name, formatter.FormatWeightFloat(r.Value)))
	}

	sb.WriteString(fmt.Sprintf("\nСредний 1ПМ: %s кг\n\n", formatter.FormatWeightFloat(summary.Average)))
	sb.WriteString("Проценты от 1ПМ:\n")
	percentages := []int{50, 60, 70, 75, 80, 85, 90, 95, 100}
	for _, p := range percentages {
		val := summary.Average * float64(p) / 100
		sb.WriteString(fmt.Sprintf("• %d%%: %s кг\n", p, formatter.FormatWeightFloat(val)))
	}

	a.sendMessage(peerID, sb.String(), nil)
}

func (a *API) StartOneRMHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	a.setUserState(userID, entity.StateAwaitingOneRMInput)
	a.sendMessage(peerID, startOneRMText, nil)
}

func (a *API) StartExerciseProgressionChartHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	keyboard := a.createMuscleGroupKeyboard()
	a.setUserState(userID, entity.StateAwaitingExerciseProgression)
	a.sendMessage(peerID, startProgressionMuscleGroupSelectText, keyboard)
}

func (a *API) ExerciseProgressionChartHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID
	userID := a.getUserIDFromCallback(obj)
	conversationMessageID := obj.ConversationMessageID
	payload := parseCallbackPayload(obj.Payload)
	exerciseIDStr := strings.TrimPrefix(payload, startGetExerciseProgressionPrefix)

	defer a.clearUserState(userID)

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		a.sendMessage(peerID, errInvalidExerciseID, nil)
		return
	}

	a.editMessage(peerID, conversationMessageID, loadingProgressionText, nil)

	data, err := a.trainingService.GetExerciseProgression(a.ctx, userID, exerciseID)
	if err != nil {
		a.sendMessage(peerID, errGetTrainings, nil)
		return
	}
	if len(data) == 0 {
		a.sendMessage(peerID, notFoundTrainingsText, nil)
		return
	}

	var xValues []string
	var yValues []float32

	for _, v := range data {
		xValues = append(xValues, v.SessionDate.Format(time.DateOnly))
		yValues = append(yValues, v.Weight)
	}

	exerciseName := data[0].ExerciseName
	cfg := chart.LinearChartConfig{
		Title:    exerciseName,
		XName:    "Дата",
		YName:    "Вес (кг)",
		YValues:  yValues,
		XValues:  xValues,
		FileName: fmt.Sprintf("%s/%s-%s.png", a.cfg.GraphicsPath, userID, time.Now().Format(time.DateOnly)),
	}

	if err = a.chartService.GenerateLinearChart(cfg); err != nil {
		a.sendMessage(peerID, errProgression, nil)
		return
	}

	a.sendPhoto(peerID, cfg.FileName)
}

func (a *API) StartExerciseHistoryHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	keyboard := a.createMuscleGroupKeyboard()
	a.setUserState(userID, entity.StateAwaitingExerciseHistory)
	a.sendMessage(peerID, startExerciseHistoryMuscleGroupSelectText, keyboard)
}

func (a *API) ExerciseHistoryHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID
	userID := a.getUserIDFromCallback(obj)
	conversationMessageID := obj.ConversationMessageID
	payload := parseCallbackPayload(obj.Payload)
	exerciseIDStr := strings.TrimPrefix(payload, startGetExerciseHistoryPrefix)

	defer a.clearUserState(userID)

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		a.sendMessage(peerID, errInvalidExerciseID, nil)

		return
	}

	sets, err := a.trainingService.GetLastSetsForExercise(a.ctx, userID, exerciseID, daysForExerciseHistory)
	if err != nil {
		a.sendMessage(peerID, errInternal, nil)

		return
	}

	if len(sets) == 0 {
		a.sendMessage(peerID, notFoundTrainingsText, nil)

		return
	}

	lastSets := a.formatter.FormatLastSets(sets)
	msgText := fmt.Sprintf(lastSetsText, lastSets)

	a.editMessage(peerID, conversationMessageID, msgText, nil)
}

func (a *API) StartGetTrainingsHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	a.setUserState(userID, entity.StateAwaitingGetTrainingsInput)
	a.sendMessage(peerID, startGetTrainingsText, nil)
}

func (a *API) GetTrainingsHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	defer a.clearUserState(userID)

	var fromDate, toDate *time.Time
	args := strings.SplitN(obj.Message.Text, " ", 2)

	if len(args) >= 1 {
		from, err := time.Parse(time.DateOnly, args[0])
		if err == nil {
			fromDate = &from
		}
	}
	if len(args) == 2 {
		to, err := time.Parse(time.DateOnly, args[1])
		if err == nil {
			toDate = &to
		}
	}

	trainings, err := a.trainingService.GetTrainingSessions(a.ctx, userID, fromDate, toDate)
	if err != nil {
		a.sendMessage(peerID, errGetTrainings, nil)

		return
	}
	if len(trainings) == 0 {
		a.sendMessage(peerID, notFoundTrainingsText, nil)

		return
	}

	text := a.formatter.FormatTrainingLogs(trainings)
	chunks := splitMessage(text, maxVKMessageLength)

	for _, chunk := range chunks {
		a.sendMessage(peerID, chunk, nil)
	}
}

func (a *API) UnknownCommandHandler(obj *events.MessageNewObject) {
	a.sendMessage(obj.Message.PeerID, unknownCommandText, nil)
}

func (a *API) StartCreateExerciseHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	if a.cfg.AuthorID != 0 && obj.Message.FromID != a.cfg.AuthorID {
		a.sendMessage(peerID, adminOnlyText, nil)

		return
	}

	a.setUserState(userID, entity.StateAwaitingExerciseInput)
	a.sendMessage(peerID, startCreateExerciseText, nil)
}

func (a *API) CreateExerciseHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	defer a.clearUserState(userID)

	lines := strings.Split(strings.TrimSpace(obj.Message.Text), "\n")
	if len(lines) < 3 {
		a.sendMessage(peerID, fmt.Sprintf(errGeneral, startCreateExerciseText), nil)

		return
	}

	name, muscleGroup, equipment := strings.TrimSpace(lines[0]), strings.TrimSpace(lines[1]), strings.TrimSpace(lines[2])
	if name == "" {
		a.sendMessage(peerID, fmt.Sprintf(errGeneral, emptyExerciseNameText), nil)

		return
	}
	if !slices.Contains(availableMuscleGroups, muscleGroup) {
		a.sendMessage(peerID, fmt.Sprintf(errGeneral, fmt.Sprintf(unknownMuscleGroupText, availableMuscleGroups)), nil)

		return
	}

	err := a.trainingService.CreateExercise(a.ctx, name, muscleGroup, equipment)
	if err != nil {
		if errors.Is(err, errs.ErrExerciseAlreadyExists) {
			a.sendMessage(peerID, fmt.Sprintf(errGeneral, fmt.Sprintf(exerciseWithNameAlreadyExistsText, name)), nil)
		} else {
			a.sendMessage(peerID, errCreateExercise, nil)
		}

		return
	}

	a.sendMessage(peerID, fmt.Sprintf(exerciseCreatedText, name, muscleGroup), nil)
}

func (a *API) StartUploadTrainingHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	a.setUserState(userID, entity.StateAwaitingTrainingInput)
	a.sendMessage(peerID, startUploadTrainingText, nil)
}

func (a *API) UploadTrainingHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	defer a.clearUserState(userID)

	session, err := a.trainingService.ParseTraining(a.ctx, entity.Event{UserID: userID, Text: obj.Message.Text})
	if err != nil {
		a.sendMessage(peerID, errUploadTraining, nil)

		return
	}

	text := fmt.Sprintf(finishText, session.ExerciseCount(), session.SetCount(), session.TotalVolume())
	a.sendMessage(peerID, text, nil)
}

func (a *API) ClearTrainingHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	defer a.clearUserState(userID)

	if err := a.trainingService.ClearSession(a.ctx, userID); err != nil {
		text := errClearTraining
		if errors.Is(err, errs.ErrSessionNotFound) {
			text = errNoTraining
		}
		a.sendMessage(peerID, text, nil)

		return
	}

	a.sendMessage(peerID, clearTrainingDoneText, nil)
}

func (a *API) StartTrainingHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)

	_, err := a.trainingService.StartTraining(a.ctx, userID)
	if err != nil {
		a.sendMessage(peerID, fmt.Sprintf(errStartTraining, err), nil)

		return
	}

	keyboard := a.createMuscleGroupKeyboard()
	a.sendMessage(peerID, startTrainingText, keyboard)
}

func (a *API) MuscleGroupHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID
	conversationMessageID := obj.ConversationMessageID
	userID := a.getUserIDFromCallback(obj)
	state := a.getUserState(userID)

	payload := parseCallbackPayload(obj.Payload)
	muscleGroup, page, _, cancelExerciseID, err := parseMuscleGroupCallbackData(payload)
	if err != nil {
		a.sendMessage(peerID, fmt.Sprintf(errGeneral, err), nil)

		return
	}

	if cancelExerciseID != uuid.Nil {
		if err := a.trainingService.DeleteExercise(a.ctx, userID, cancelExerciseID); err != nil {
			a.sendMessage(peerID, fmt.Sprintf(errGeneral, err), nil)

			return
		}
	}

	exercises, err := a.trainingService.GetExercisesByMuscleGroup(a.ctx, muscleGroup)
	if err != nil {
		a.sendMessage(peerID, errExerciseLoad, nil)

		return
	}
	if len(exercises) == 0 {
		a.sendMessage(peerID, errNoExercises, nil)

		return
	}

	pagedExercises, totalPages, err := paginate(exercises, page, pageSize)
	if err != nil {
		a.sendMessage(peerID, err.Error(), nil)

		return
	}

	callbackDataPrefix := fmt.Sprintf("%s%s:", exercisePrefix, muscleGroup)
	switch state {
	case entity.StateAwaitingExerciseProgression:
		callbackDataPrefix = startGetExerciseProgressionPrefix
	case entity.StateAwaitingExerciseHistory:
		callbackDataPrefix = startGetExerciseHistoryPrefix
	}

	keyboard := a.createExerciseKeyboard(pagedExercises, muscleGroup, page, totalPages, callbackDataPrefix)
	a.editMessage(peerID, conversationMessageID, fmt.Sprintf(muscleGroupDoneText, muscleGroup), keyboard)
}

func (a *API) BackToMuscleGroupsHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID
	conversationMessageID := obj.ConversationMessageID

	keyboard := a.createMuscleGroupKeyboard()
	a.editMessage(peerID, conversationMessageID, muscleGroupSelectText, keyboard)
}

func (a *API) ExerciseHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID
	conversationMessageID := obj.ConversationMessageID
	userID := a.getUserIDFromCallback(obj)
	payload := parseCallbackPayload(obj.Payload)
	text := strings.TrimPrefix(payload, exercisePrefix)
	args := strings.SplitN(text, ":", 2)

	if len(args) != 2 {
		a.sendMessage(peerID, errInternal, nil)

		return
	}

	muscleGroup, exerciseIDStr := args[0], args[1]

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		a.sendMessage(peerID, errInvalidExerciseID, nil)

		return
	}

	err = a.trainingService.AddTrainingExercise(a.ctx, userID, exerciseID)
	if err != nil {
		a.sendMessage(peerID, fmt.Sprintf(errAddExercise, err), nil)

		return
	}

	sets, err := a.trainingService.GetLastSetsForExercise(a.ctx, userID, exerciseID, daysForSetStatistics)
	if err != nil {
		a.sendMessage(peerID, errInternal, nil)

		return
	}

	lastSets := a.formatter.FormatLastSets(sets)

	keyboard := a.createBackToExercisesKeyboard(muscleGroup, exerciseIDStr)

	msgText := exerciseText
	if lastSets != "" {
		msgText = fmt.Sprintf("%s\n\n%s", exerciseText, fmt.Sprintf(lastSetsText, lastSets))
	}

	a.setUserState(userID, entity.StateAwaitingSetInput)
	a.editMessage(peerID, conversationMessageID, msgText, keyboard)
}

func (a *API) SetHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := a.getUserID(obj)
	input := obj.Message.Text

	parts := strings.SplitN(input, "\n", 2)
	setData := strings.Split(parts[0], ",")
	if len(setData) != 2 {
		a.sendMessage(peerID, errInvalidFormat, nil)

		return
	}

	weight, errWeight := strconv.ParseFloat(setData[0], 64)
	reps, errReps := strconv.Atoi(setData[1])
	if errWeight != nil || errReps != nil {
		a.sendMessage(peerID, errParseData, nil)

		return
	}
	var notes string
	if len(parts) > 1 {
		notes = strings.TrimSpace(parts[1])
	}

	err := a.trainingService.AddOrUpdateSet(a.ctx, userID, obj.Message.ID, float32(weight), uint8(reps), notes)
	if err != nil {
		a.sendMessage(peerID, fmt.Sprintf(errGeneral, err), nil)

		return
	}

	keyboard := a.createSetActionsKeyboard()
	a.sendMessage(peerID, setText, keyboard)
}

func (a *API) StartNewExerciseHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID
	userID := a.getUserIDFromCallback(obj)

	defer a.clearUserState(userID)

	keyboard := a.createMuscleGroupKeyboard()
	a.sendMessage(peerID, muscleGroupSelectText, keyboard)
}

func (a *API) ConfirmationFinishTrainingHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID

	keyboard := a.createFinishConfirmationKeyboard()
	a.sendMessage(peerID, finishTrainingConfirmationText, keyboard)
}

func (a *API) AcceptFinishTrainingHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID
	conversationMessageID := obj.ConversationMessageID
	userID := a.getUserIDFromCallback(obj)

	defer a.clearUserState(userID)

	session, err := a.trainingService.EndSession(a.ctx, userID)
	if err != nil {
		a.sendMessage(peerID, fmt.Sprintf(errGeneral, err), nil)

		return
	}

	text := fmt.Sprintf(finishText, session.ExerciseCount(), session.SetCount(), session.TotalVolume())
	a.editMessage(peerID, conversationMessageID, text, nil)

	details := a.formatter.FormatTrainingLogs([]entity.TrainingSession{*session})
	if details != "" {
		chunks := splitMessage(details, maxVKMessageLength)
		for _, chunk := range chunks {
			a.sendMessage(peerID, chunk, nil)
		}
	}
}

func (a *API) RejectFinishTrainingHandler(obj *events.MessageEventObject) {
	a.sendEventAnswer(obj.EventID, obj.UserID, obj.PeerID)

	peerID := obj.PeerID
	conversationMessageID := obj.ConversationMessageID

	a.editMessage(peerID, conversationMessageID, "Продолжайте тренировку!", nil)
}

func parseMuscleGroupCallbackData(data string) (muscleGroup string, page int, direction string, exerciseID uuid.UUID, err error) {
	parts := strings.Split(data, ":")

	if len(parts) < 2 {
		return "", 0, "", uuid.Nil, fmt.Errorf("invalid callback format")
	}

	muscleGroup = parts[1]
	page = 0
	direction = ""

	if len(parts) == 4 {
		direction = parts[2]
		page, err = strconv.Atoi(parts[3])
		if err != nil {
			return "", 0, "", uuid.Nil, fmt.Errorf("invalid page number")
		}
	}

	if len(parts) == 5 {
		exerciseIDStr := parts[4]
		exerciseID, err = uuid.Parse(exerciseIDStr)
		if err != nil {
			return "", 0, "", uuid.Nil, fmt.Errorf("invalid exercise id")
		}
	}

	return muscleGroup, page, direction, exerciseID, nil
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

func splitMessage(text string, maxLen int) []string {
	var chunks []string

	for len(text) > maxLen {
		splitAt := maxLen

		if idx := strings.LastIndex(text[:splitAt], "\n"); idx > 0 {
			splitAt = idx
		}

		chunks = append(chunks, text[:splitAt])
		text = text[splitAt:]
	}

	if len(text) > 0 {
		chunks = append(chunks, text)
	}

	return chunks
}
