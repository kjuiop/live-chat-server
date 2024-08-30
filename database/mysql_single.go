package database

import (
	"context"
	"live-chat-server/config"
	"time"
)

type mysqlClient struct {
	cfg config.Mysql
}

func (m mysqlClient) HSet(ctx context.Context, key string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m mysqlClient) Expire(ctx context.Context, key string, expTime time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (m mysqlClient) HGet(ctx context.Context, key, mapKey string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m mysqlClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m mysqlClient) Exists(ctx context.Context, key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m mysqlClient) DelByKey(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func NewMysqlSingleClient(ctx context.Context, cfg config.Mysql) (Client, error) {
	return &mysqlClient{
		cfg: cfg,
	}, nil
}
