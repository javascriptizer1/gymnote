package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"gymnote/internal/entity"
)

const KEY_PREFIX = "training:"

func (r *cache) SaveSession(ctx context.Context, session *entity.TrainingSession) error {
	data, err := json.Marshal(NewTrainingSessionRow(session))
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, KEY_PREFIX+session.UserID(), data, 0).Err()
}

func (r *cache) GetSession(ctx context.Context, userID string) (*entity.TrainingSession, error) {
	data, err := r.redisClient.Get(ctx, KEY_PREFIX+userID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session TrainingSessionRow
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return nil, err
	}

	return session.ToEntity(), nil
}

func (r *cache) DeleteSession(ctx context.Context, userID string) error {
	return r.redisClient.Del(ctx, KEY_PREFIX+userID).Err()
}
