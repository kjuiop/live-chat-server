package redis

import (
	"context"
	"time"
)

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}

func (m memoryClient) HSet(ctx context.Context, key string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) Expire(ctx context.Context, key string, expTime time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) Exists(ctx context.Context, key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) DelByKey(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) HGet(ctx context.Context, key, mapKey string) (string, error) {
	//TODO implement me
	panic("implement me")
}
