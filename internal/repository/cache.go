package repository

import (
	"context"

	"gymnote/internal/entity"
)

type Cache interface {
	SaveSession(ctx context.Context, session *entity.TrainingSession) error
	GetSession(ctx context.Context, userID string) (*entity.TrainingSession, error)
	DeleteSession(ctx context.Context, userID string) error
}
