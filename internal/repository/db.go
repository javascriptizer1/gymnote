package repository

import (
	"context"

	"github.com/google/uuid"

	"gymnote/internal/entity"
)

type DB interface {
	Close()

	GetExerciseByName(ctx context.Context, req string) (entity.Exercise, error)
	GetExerciseByID(ctx context.Context, req uuid.UUID) (entity.Exercise, error)
	GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error)

	InsertTrainingLogs(ctx context.Context, req entity.TrainingSession) error
	InsertTrainingSession(ctx context.Context, req entity.TrainingSession) error
}
