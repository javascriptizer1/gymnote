package repository

import (
	"context"

	"github.com/google/uuid"

	"gymnote/internal/entity"
)

type Cache interface {
	Close(_ context.Context) error
	SaveSession(ctx context.Context, session *entity.TrainingSession) error
	GetSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	DeleteSession(ctx context.Context, userID string) error
	GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error)
	SaveExercisesByMuscleGroup(ctx context.Context, muscleGroup string, exercises []entity.Exercise) error
	AddExerciseToCache(ctx context.Context, exercise entity.Exercise) error
	GetExerciseByID(ctx context.Context, exerciseID uuid.UUID) (entity.Exercise, error)
	SaveExerciseByID(ctx context.Context, exercise entity.Exercise) error
}
