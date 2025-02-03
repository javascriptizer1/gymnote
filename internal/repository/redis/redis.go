package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"gymnote/internal/config"
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
