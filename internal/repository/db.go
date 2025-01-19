package repository

import (
	"context"

	"gymnote/internal/entity"
)

type DB interface {
	Close()

	GetExerciseByName(ctx context.Context, req string) (entity.Exercise, error)

	InsertTrainingLogs(ctx context.Context, req entity.TrainingSession) error
	InsertTrainingSession(ctx context.Context, req entity.TrainingSession) error
}
