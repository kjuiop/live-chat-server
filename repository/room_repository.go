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

	if err := r.rdb.HSet(ctx, room.GenerateRoomKey(), room.ConvertRedisData()); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	if err := r.rdb.Expire(ctx, room.RoomId, time.Duration(room.RoomIdTTLDay*24)*time.Hour); err != nil {
		return fmt.Errorf("create chat room key fail set ttl, key : %s, err : %w", room.RoomId, err)
	}

	return nil

}

func (r *roomRepository) Fetch(ctx context.Context, key string) (models.RoomInfo, error) {

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

func (r *roomRepository) Exists(ctx context.Context, key string) (bool, error) {

	isExist, err := r.rdb.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("fail redis cmd exist err : %w", err)
	}

	return isExist, nil
}

func (r *roomRepository) Update(ctx context.Context, key string, room *models.RoomInfo) error {

	if err := r.rdb.HSet(ctx, room.RoomId, room.ConvertRedisData()); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	if err := r.rdb.Expire(ctx, room.RoomId, time.Duration(room.RoomIdTTLDay*24)*time.Hour); err != nil {
		return fmt.Errorf("create chat room key fail set ttl, key : %s, err : %w", room.RoomId, err)
	}

	return nil
}

func (r *roomRepository) Delete(ctx context.Context, key string) error {

	if err := r.rdb.DelByKey(ctx, key); err != nil {
		return err
	}

	return nil
}

func (r *roomRepository) SetRoomMap(ctx context.Context, key string, data *models.RoomInfo) error {

	jData, err := json.Marshal(data.ConvertRedisData())
	if err != nil {
		return fmt.Errorf("set room map json encoding fail, err : %w", err)
	}

	if err := r.rdb.HSet(ctx, key, map[string]interface{}{
		data.GenerateRoomMapFieldKey(): string(jData),
	}); err != nil {
		return err
	}

	return nil
}

func (r *roomRepository) GetRoomMap(ctx context.Context, key, mapKey string) (models.RoomInfo, error) {

	result, err := r.rdb.HGet(ctx, key, mapKey)
	if err != nil {
		return models.RoomInfo{}, err
	}

	var roomInfo models.RoomInfo
	if err := json.Unmarshal([]byte(result), &roomInfo); err != nil {
		return models.RoomInfo{}, fmt.Errorf("json parsing err : %w", err)
	}

	return roomInfo, nil
}
