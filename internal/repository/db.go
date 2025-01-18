package repository

import (
	"context"

	"gymnote/internal/entity"

	"github.com/google/uuid"
)

type DB interface {
	Close()

	GetExerciseIDByName(ctx context.Context, req string) (uuid.UUID, error)

	InsertTrainingLogs(ctx context.Context, req entity.TrainingSession) error
	InsertTrainingSession(ctx context.Context, req entity.TrainingSession) error
}
