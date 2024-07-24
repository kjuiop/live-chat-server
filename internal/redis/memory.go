package repository

import (
	"context"
	"time"
)

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}

func (m memoryClient) HMSet(ctx context.Context, key string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) Expire(ctx context.Context, key string, expTime time.Duration) error {
	//TODO implement me
	panic("implement me")
}
