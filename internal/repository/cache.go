package repository

import (
	"context"

	"gymnote/internal/entity"
)

type Cache interface {
	Close(_ context.Context) error
	SaveSession(ctx context.Context, session *entity.TrainingSession) error
	GetSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	DeleteSession(ctx context.Context, userID string) error
}
