package repository

import (
	"context"
	"fmt"
	redis "live-chat-server/internal/redis"
	"live-chat-server/models"
	"time"
)

type roomRepository struct {
	rdb redis.Client
}

func NewRoomRepository(client redis.Client) models.RoomRepository {
	return &roomRepository{
		rdb: client,
	}
}

func (r roomRepository) Create(ctx context.Context, room *models.RoomInfo) error {

	if room.RoomId == "" {
		return fmt.Errorf("required chat room id : %s", room.RoomId)
	}

	if err := r.rdb.HMSet(ctx, room.RoomId, room.ConvertRedisData()); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	if err := r.rdb.Expire(ctx, room.RoomId, time.Duration(2)*time.Hour); err != nil {
		return fmt.Errorf("create chat room key fail set ttl, key : %s, err : %w", room.RoomId, err)
	}

	return nil

}
