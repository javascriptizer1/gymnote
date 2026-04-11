package vk

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/google/uuid"

	"gymnote/internal/chart"
	"gymnote/internal/config"
	"gymnote/internal/entity"
)

type CommandHandler func(*events.MessageNewObject)
type CallbackHandler func(*events.MessageEventObject)

type Formatter interface {
	FormatTrainingLogs(sessions []entity.TrainingSession) string
	FormatLastSets(sessions []entity.ExerciseProgression) string
}

type ChartService interface {
	GenerateLinearChart(config chart.LinearChartConfig) error
}

type AccountService interface {
	GetOrCreateAccount(ctx context.Context, platform entity.Platform, platformUserID string, username string) (*entity.UserAccount, error)
	CreateLinkCode(ctx context.Context, platform entity.Platform, platformUserID string) (*entity.LinkCode, error)
	LinkAccountWithCode(ctx context.Context, code string, platform entity.Platform, platformUserID string, username string) error
	GetLinkedAccounts(ctx context.Context, platform entity.Platform, platformUserID string) ([]entity.UserAccount, error)
}

type TrainingService interface {
	ParseTraining(ctx context.Context, e entity.Event) (*entity.TrainingSession, error)
	GetExerciseProgression(ctx context.Context, userID string, exerciseID uuid.UUID) ([]entity.ExerciseProgression, error)
	GetTrainingSessions(ctx context.Context, userID string, fromDate, toDate *time.Time) ([]entity.TrainingSession, error)
	GetLastSetsForExercise(ctx context.Context, userID string, exerciseID uuid.UUID, limitDays int64) ([]entity.ExerciseProgression, error)
	DeleteExercise(ctx context.Context, userID string, exerciseID uuid.UUID) error
	CreateExercise(ctx context.Context, name string, muscleGroup string, equipment string) error
	StartTraining(ctx context.Context, userID string) (*entity.TrainingSession, error)
	AddTrainingExercise(ctx context.Context, userID string, exerciseID uuid.UUID) error
	AddOrUpdateSet(ctx context.Context, userID string, messageID int, weight float32, reps uint8, notes string) error
	EndSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	GetCurrentSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	ClearSession(ctx context.Context, userID string) error
	GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error)
	UpdateSetFromMessage(ctx context.Context, userID string, messageID int, weight float32, reps uint8, notes string) error
}

type API struct {
	ctx              context.Context
	cfg              *config.VKConfig
	vk               *api.VK
	formatter        Formatter
	chartService     ChartService
	trainingService  TrainingService
	accountService   AccountService
	commandHandlers  map[string]CommandHandler
	stateHandlers    map[entity.UserState]func(*events.MessageNewObject)
	callbackHandlers map[string]CallbackHandler
	userStates       map[string]entity.UserState
	mu               sync.Mutex
}

func NewAPI(ctx context.Context, cfg *config.VKConfig, formatter Formatter, chartService ChartService, trainingService TrainingService, accountService AccountService) *API {
	vk := api.NewVK(cfg.AccessToken)

	api := &API{
		ctx:              ctx,
		cfg:              cfg,
		vk:               vk,
		formatter:        formatter,
		chartService:     chartService,
		trainingService:  trainingService,
		accountService:   accountService,
		commandHandlers:  make(map[string]CommandHandler),
		callbackHandlers: make(map[string]CallbackHandler),
		stateHandlers:    make(map[entity.UserState]func(*events.MessageNewObject)),
		userStates:       make(map[string]entity.UserState),
		mu:               sync.Mutex{},
	}

	api.registerHandlers()
	api.setBotCommands()

	log.Println("VK Bot API initialized")

	return api
}

func (a *API) registerHandlers() {
	a.commandHandlers = map[string]CommandHandler{
		startCommand:                  a.StartHandler,
		helpCommand:                   a.HelpHandler,
		startTrainingCommand:          a.StartTrainingHandler,
		createExerciseCommand:         a.StartCreateExerciseHandler,
		clearTrainingCommand:          a.ClearTrainingHandler,
		uploadTrainingCommand:         a.StartUploadTrainingHandler,
		getTrainingsCommand:           a.StartGetTrainingsHandler,
		getExerciseProgressionCommand: a.StartExerciseProgressionChartHandler,
		getExerciseHistoryCommand:     a.StartExerciseHistoryHandler,
		oneRMCommand:                  a.StartOneRMHandler,
		linkAccountCommand:            a.LinkAccountHandler,
		myAccountsCommand:             a.MyAccountsHandler,

		"🏋️ Начать тренировку":   a.StartTrainingHandler,
		"📝 Загрузить тренировку": a.StartUploadTrainingHandler,
		"📊 История тренировок":   a.StartGetTrainingsHandler,
		"📈 Прогрессия весов":     a.StartExerciseProgressionChartHandler,
		"📋 История упражнения":   a.StartExerciseHistoryHandler,
		"🔢 Расчет 1ПМ":           a.StartOneRMHandler,
		"🗑 Сбросить тренировку":  a.ClearTrainingHandler,
		"🔗 Связать аккаунты":     a.LinkAccountHandler,
		"👤 Мои аккаунты":         a.MyAccountsHandler,
		"➕ Создать упражнение":   a.StartCreateExerciseHandler,
		"❓ Помощь":               a.HelpHandler,
	}

	a.stateHandlers = map[entity.UserState]func(*events.MessageNewObject){
		entity.StateAwaitingSetInput:          a.SetHandler,
		entity.StateAwaitingExerciseInput:     a.CreateExerciseHandler,
		entity.StateAwaitingTrainingInput:     a.UploadTrainingHandler,
		entity.StateAwaitingGetTrainingsInput: a.GetTrainingsHandler,
		entity.StateAwaitingOneRMInput:        a.OneRMHandler,
	}

	a.callbackHandlers = map[string]CallbackHandler{
		musclePrefix:                      a.MuscleGroupHandler,
		exercisePrefix:                    a.ExerciseHandler,
		confirmationFinishTrainingPrefix:  a.ConfirmationFinishTrainingHandler,
		acceptFinishTrainingPrefix:        a.AcceptFinishTrainingHandler,
		rejectFinishTrainingPrefix:        a.RejectFinishTrainingHandler,
		startNewExercisePrefix:            a.StartNewExerciseHandler,
		startGetExerciseProgressionPrefix: a.ExerciseProgressionChartHandler,
		startGetExerciseHistoryPrefix:     a.ExerciseHistoryHandler,
		backToMuscleGroups:                a.BackToMuscleGroupsHandler,
	}
}

func (a *API) setBotCommands() {
	type Command struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	commands := []Command{
		{Name: startCommand, Description: "Запустить бота"},
		{Name: startTrainingCommand, Description: "Начать тренировку"},
		{Name: uploadTrainingCommand, Description: "Загрузить тренировку"},
		{Name: getTrainingsCommand, Description: "Посмотреть историю тренировок"},
		{Name: getExerciseProgressionCommand, Description: "Посмотреть прогрессию весов по упражнению"},
		{Name: getExerciseHistoryCommand, Description: "Посмотреть историю конкретного упражнения"},
		{Name: createExerciseCommand, Description: "Создать новое упражнение"},
		{Name: clearTrainingCommand, Description: "Сбросить текущую тренировку"},
		{Name: oneRMCommand, Description: "Рассчитать одноповторный максимум"},
		{Name: linkAccountCommand, Description: "Связать аккаунты Telegram и VK"},
		{Name: myAccountsCommand, Description: "Мои связанные аккаунты"},
		{Name: helpCommand, Description: "Помощь и команды"},
	}

	if _, err := a.vk.Request("messages.setChatMenuButton", api.Params{
		"group_id": a.cfg.GroupID,
		"type":     "commands",
		"commands": commands,
	}); err != nil {
		log.Printf("VK: Failed to set bot commands: %v", err)
	} else {
		log.Printf("VK: Bot commands configured")
	}
}
