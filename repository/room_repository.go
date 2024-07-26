package repository

import (
	"context"
	"encoding/json"
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

func (r *roomRepository) Create(ctx context.Context, room *models.RoomInfo) error {

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

func (r *roomRepository) Fetch(ctx context.Context, key string) (models.RoomInfo, error) {

	if key == "" {
		return models.RoomInfo{}, fmt.Errorf("fail get room info, required chat room id : %s", key)
	}

	result, err := r.rdb.HGetAll(ctx, key)
	if err != nil {
		return models.RoomInfo{}, fmt.Errorf("fail get room info, err : %w", err)
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return models.RoomInfo{}, fmt.Errorf("json marshal err : %w", err)
	}

	var roomInfo models.RoomInfo
	if err := json.Unmarshal(jsonData, &roomInfo); err != nil {
		return models.RoomInfo{}, fmt.Errorf("json parsing err : %w", err)
	}

	return roomInfo, nil
}
