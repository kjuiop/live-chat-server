package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"live-chat-server/config"
	"live-chat-server/models"
	"time"
)

type redisClient struct {
	cfg    config.Redis
	client *redis.Client
}

func NewRedisSingleClient(ctx context.Context, cfg config.Redis) (Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("fail ping err : %w", err)
	}

	return &redisClient{
		cfg:    cfg,
		client: client,
	}, nil
}

func (r *redisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {

	if key == "" {
		return nil, errors.New(models.EmptyRedisKey)
	}

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("fail get data hgetall, err : %w", err)
	}

	return result, nil
}

func (r *redisClient) HMSet(ctx context.Context, key string, data map[string]interface{}) error {

	if key == "" {
		return errors.New(models.EmptyRedisKey)
	}

	if err := r.client.HMSet(ctx, key, data).Err(); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	return nil
}

func (r *redisClient) Expire(ctx context.Context, key string, expTime time.Duration) error {

	if key == "" {
		return errors.New(models.EmptyRedisKey)
	}

	if err := r.client.Expire(ctx, key, expTime).Err(); err != nil {
		return fmt.Errorf("fail set ttl, key : %w", err)
	}

	return nil
}
