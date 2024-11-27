package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"live-chat-server/config"
	"strings"
	"time"
)

type universal struct {
	cfg  config.Redis
	conn redis.UniversalClient
}

func NewUniversalClient(ctx context.Context, cfg config.Redis) (Client, error) {

	client := &universal{
		cfg: cfg,
	}

	switch cfg.Mode {
	case "single":
		client.conn = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        strings.Split(cfg.Addr, ","),
			PoolSize:     cfg.PoolSize,
			DialTimeout:  time.Second * 3,
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
		})
	case "sentinel":
		client.conn = redis.NewUniversalClient(&redis.UniversalOptions{
			MasterName:   cfg.Masters,
			Addrs:        strings.Split(cfg.Addr, ","),
			PoolSize:     cfg.PoolSize,
			DialTimeout:  time.Second * 3,
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
		})
	case "cluster":
		client.conn = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        strings.Split(cfg.Addr, ","),
			PoolSize:     cfg.PoolSize,
			Password:     cfg.Password,
			DialTimeout:  time.Second * 3,
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
		})
	default:
		return nil, errors.New("invalid redis mode")
	}

	if err := client.conn.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func (u universal) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (u universal) Get(ctx context.Context, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (u universal) HSet(ctx context.Context, key, fieldKey string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (u universal) Expire(ctx context.Context, key string, expTime time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (u universal) HGet(ctx context.Context, key, mapKey string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (u universal) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (u universal) Exists(ctx context.Context, key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (u universal) DelByKey(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}
