package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"gymnote/internal/entity"
)

type TrainingService interface {
	StartTraining(ctx context.Context, userID string) (*entity.TrainingSession, error)
	AddExercise(ctx context.Context, userID string, exerciseID uuid.UUID) error
	AddSet(ctx context.Context, userID string, set *entity.Set) error
	UpdateActiveSet(ctx context.Context, userID string, weight float32, reps uint8) error
	EndSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	GetCurrentSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error)
}

type API struct {
	ctx             context.Context
	bot             *tgbotapi.BotAPI
	trainingService TrainingService
}

func New(ctx context.Context, bot *tgbotapi.BotAPI, trainingService TrainingService) *API {
	return &API{
		ctx:             ctx,
		bot:             bot,
		trainingService: trainingService,
	}
}
