package tg

import (
	"context"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"gymnote/internal/chart"
	"gymnote/internal/config"
	"gymnote/internal/entity"
)

type CommandHandler func(*tgbotapi.Message)
type CallbackHandler func(*tgbotapi.CallbackQuery)

type Formatter interface {
	FormatTrainingLogs(sessions []entity.TrainingSession) string
	FormatLastSets(sessions []entity.ExerciseProgression) string
}
type ChartService interface {
	GenerateLinearChart(config chart.LinearChartConfig) error
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
	AddOrUpdateSet(ctx context.Context, userID string, weight float32, reps uint8, notes string) error
	EndSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	GetCurrentSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	ClearSession(ctx context.Context, userID string) error
	GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error)
}

type API struct {
	ctx              context.Context
	cfg              *config.TelegramConfig
	bot              *tgbotapi.BotAPI
	formatter        Formatter
	chartService     ChartService
	trainingService  TrainingService
	commandHandlers  map[string]CommandHandler
	stateHandlers    map[entity.UserState]func(*tgbotapi.Message)
	callbackHandlers map[string]CallbackHandler
	userStates       map[string]entity.UserState
	mu               sync.Mutex
}

func NewAPI(ctx context.Context, cfg *config.TelegramConfig, formatter Formatter, chartService ChartService, trainingService TrainingService) *API {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = cfg.Debug

	api := &API{
		ctx:              ctx,
		cfg:              cfg,
		bot:              bot,
		formatter:        formatter,
		chartService:     chartService,
		trainingService:  trainingService,
		commandHandlers:  make(map[string]CommandHandler),
		callbackHandlers: make(map[string]CallbackHandler),
		stateHandlers:    make(map[entity.UserState]func(*tgbotapi.Message)),
		userStates:       make(map[string]entity.UserState),
		mu:               sync.Mutex{},
	}

	api.setBotCommands()
	api.registerHandlers()

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
	}

	a.stateHandlers = map[entity.UserState]func(*tgbotapi.Message){
		entity.StateAwaitingSetInput:          a.SetHandler,
		entity.StateAwaitingExerciseInput:     a.CreateExerciseHandler,
		entity.StateAwaitingTrainingInput:     a.UploadTrainingHandler,
		entity.StateAwaitingGetTrainingsInput: a.GetTrainingsHandler,
	}

	a.callbackHandlers = map[string]CallbackHandler{
		musclePrefix:                      a.MuscleGroupHandler,
		exercisePrefix:                    a.ExerciseHandler,
		finishTrainingPrefix:              a.FinishTrainingHandler,
		startNewExercisePrefix:            a.StartNewExerciseHandler,
		startGetExerciseProgressionPrefix: a.ExerciseProgressionChartHandler,
		backToMuscleGroups:                a.BackToMuscleGroupsHandler,
	}
}

func (a *API) setBotCommands() {
	commands := []tgbotapi.BotCommand{
		{Command: startCommand, Description: "Запустить бота"},
		{Command: startTrainingCommand, Description: "Начать тренировку"},
		{Command: uploadTrainingCommand, Description: "Загрузить тренировку"},
		{Command: getTrainingsCommand, Description: "Посмотреть историю тренировок"},
		{Command: getExerciseProgressionCommand, Description: "Посмотреть прогрессию весов по упражнению"},
		{Command: createExerciseCommand, Description: "Создать новое упражнение"},
		{Command: clearTrainingCommand, Description: "Сбросить текущую тренировку"},
		{Command: helpCommand, Description: "Помощь и команды"},
	}

	_, err := a.bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Printf("Set commands error %v", err)
	}
}
