package repository

import (
	"context"
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

func NewRedisClient(ctx context.Context, cfg config.Redis) (Client, error) {

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

func (r redisClient) CreateChatRoom(ctx context.Context, room *models.RoomInfo) error {

	if room.RoomId == "" {
		return fmt.Errorf("required chat room id : %s", room.RoomId)
	}

	if err := r.client.HMSet(ctx, room.RoomId, room.ConvertRedisData()).Err(); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	if err := r.client.Expire(ctx, room.RoomId, time.Duration(2)*time.Hour).Err(); err != nil {
		return fmt.Errorf("create chat room key fail set ttl, key : %s, err : %w", room.RoomId, err)
	}

	return nil
}
