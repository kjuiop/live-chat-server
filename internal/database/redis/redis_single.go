package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"live-chat-server/config"
	"time"
)

// 쿼리가 실행되는 공간

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

func (r *redisClient) HGet(ctx context.Context, key, mapKey string) (string, error) {

	result, err := r.client.HGet(ctx, key, mapKey).Result()
	if err != nil {
		return "", fmt.Errorf("fail hget data, err : %w", err)
	}

	return result, nil
}

func (r *redisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("fail get data hgetall, err : %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("not exist room key : %s", key)
	}

	return result, nil
}

func (r *redisClient) HSet(ctx context.Context, key string, data map[string]interface{}) error {

	if err := r.client.HSet(ctx, key, data).Err(); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	return nil
}

func (r *redisClient) Expire(ctx context.Context, key string, expTime time.Duration) error {

	if err := r.client.Expire(ctx, key, expTime).Err(); err != nil {
		return fmt.Errorf("fail set ttl, key : %w", err)
	}

	return nil
}

func (r *redisClient) Exists(ctx context.Context, key string) (bool, error) {
	isExist, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return isExist == 1, nil
}

func (r *redisClient) DelByKey(ctx context.Context, key string) error {

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
