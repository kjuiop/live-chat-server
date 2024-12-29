package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"live-chat-server/config"
	"time"
)

type sentinelClient struct {
	cfg    config.Redis
	client *redis.Client
}

func NewRedisSentinelClient(ctx context.Context, cfg config.Redis) (Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("fail ping err : %w", err)
	}

	return &sentinelClient{
		cfg:    cfg,
		client: client,
	}, nil
}

func (s sentinelClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (s sentinelClient) Get(ctx context.Context, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s sentinelClient) HSet(ctx context.Context, key, fieldKey string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s sentinelClient) Expire(ctx context.Context, key string, expTime time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (s sentinelClient) HGet(ctx context.Context, key, mapKey string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s sentinelClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s sentinelClient) Exists(ctx context.Context, key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s sentinelClient) DelByKey(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}
