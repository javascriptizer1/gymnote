package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"gymnote/internal/config"
	"gymnote/internal/entity"
)

type cache struct {
	redisClient *redis.Client
	cfg         *config.CacheConfig
}

func New(ctx context.Context, cfg *config.CacheConfig) (*cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %v", err)
	}

	if pong != "PONG" {
		return nil, fmt.Errorf("unexpected response from redis: %s", pong)
	}

	return &cache{redisClient: client, cfg: cfg}, nil
}

func (c *cache) Close(_ context.Context) error {
	return c.redisClient.Close()
}

func (c *cache) exerciseKey(muscleGroup string) string {
	return fmt.Sprintf("exercises:%s", muscleGroup)
}

func (c *cache) GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error) {
	data, err := c.redisClient.Get(ctx, c.exerciseKey(muscleGroup)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var exercises []ExerciseCacheRow
	if err := json.Unmarshal(data, &exercises); err != nil {
		return nil, err
	}

	result := make([]entity.Exercise, 0, len(exercises))
	for _, e := range exercises {
		result = append(result, *e.ToEntity())
	}

	return result, nil
}

func (c *cache) SaveExercisesByMuscleGroup(ctx context.Context, muscleGroup string, exercises []entity.Exercise) error {
	rows := make([]ExerciseCacheRow, 0, len(exercises))
	for _, e := range exercises {
		rows = append(rows, NewExerciseCacheRow(&e))
	}

	data, err := json.Marshal(rows)
	if err != nil {
		return err
	}

	return c.redisClient.Set(ctx, c.exerciseKey(muscleGroup), data, 0).Err()
}

func (c *cache) AddExerciseToCache(ctx context.Context, exercise entity.Exercise) error {
	muscleGroup := exercise.MuscleGroup()

	exercises, err := c.GetExercisesByMuscleGroup(ctx, muscleGroup)
	if err != nil {
		return err
	}

	if exercises == nil {
		exercises = []entity.Exercise{}
	}

	exercises = append(exercises, exercise)

	if err := c.SaveExercisesByMuscleGroup(ctx, muscleGroup, exercises); err != nil {
		return err
	}

	return c.SaveExerciseByID(ctx, exercise)
}

func (c *cache) exerciseIDKey(exerciseID uuid.UUID) string {
	return fmt.Sprintf("exercise:%s", exerciseID.String())
}

func (c *cache) GetExerciseByID(ctx context.Context, exerciseID uuid.UUID) (entity.Exercise, error) {
	data, err := c.redisClient.Get(ctx, c.exerciseIDKey(exerciseID)).Bytes()
	if err == redis.Nil {
		return entity.Exercise{}, nil
	}
	if err != nil {
		return entity.Exercise{}, err
	}

	var row ExerciseCacheRow
	if err := json.Unmarshal(data, &row); err != nil {
		return entity.Exercise{}, err
	}

	return *row.ToEntity(), nil
}

func (c *cache) SaveExerciseByID(ctx context.Context, exercise entity.Exercise) error {
	row := NewExerciseCacheRow(&exercise)

	data, err := json.Marshal(row)
	if err != nil {
		return err
	}

	return c.redisClient.Set(ctx, c.exerciseIDKey(exercise.ID()), data, 0).Err()
}
