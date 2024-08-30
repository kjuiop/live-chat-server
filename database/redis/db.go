package redis

import (
	"context"
	"time"
)

type Client interface {
	HSet(ctx context.Context, key string, data map[string]interface{}) error
	Expire(ctx context.Context, key string, expTime time.Duration) error
	HGet(ctx context.Context, key, mapKey string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Exists(ctx context.Context, key string) (bool, error)
	DelByKey(ctx context.Context, key string) error
}
