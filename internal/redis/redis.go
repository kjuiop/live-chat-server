package repository

import (
	"context"
	"time"
)

type Client interface {
	HMSet(ctx context.Context, key string, data map[string]interface{}) error
	Expire(ctx context.Context, key string, expTime time.Duration) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
}
