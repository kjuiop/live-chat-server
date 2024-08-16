package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"live-chat-server/domain/room"
	redis "live-chat-server/internal/redis"
	"strconv"
	"time"
)

type roomRepository struct {
	rdb redis.Client
}

func NewRoomRepository(client redis.Client) room.RoomRepository {
	return &roomRepository{
		rdb: client,
	}
}

func (r *roomRepository) Create(ctx context.Context, room *room.RoomInfo) error {

	if err := r.rdb.HSet(ctx, room.GenerateRoomKey(), room.ConvertRedisData()); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	if err := r.rdb.Expire(ctx, room.RoomId, time.Duration(room.RoomIdTTLDay*24)*time.Hour); err != nil {
		return fmt.Errorf("create chat room key fail set ttl, key : %s, err : %w", room.RoomId, err)
	}

	return nil

}

func (r *roomRepository) Fetch(ctx context.Context, key string) (room.RoomInfo, error) {

	result, err := r.rdb.HGetAll(ctx, key)
	if err != nil {
		return room.RoomInfo{}, fmt.Errorf("fail get room info, err : %w", err)
	}

	// redis 는 하나의 json 문자열로 반환하기 때문에 안의 데이터타입을 변경하기 위해서는 아래와 같은 매핑 작업을 필요로 함
	data := make(map[string]interface{})
	for k, v := range result {
		if k == "created_at" {
			createdAt, _ := strconv.Atoi(v)
			data[k] = createdAt
			continue
		} else {
			data[k] = v
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return room.RoomInfo{}, fmt.Errorf("json marshal err : %w", err)
	}

	var roomInfo room.RoomInfo
	if err := json.Unmarshal(jsonData, &roomInfo); err != nil {
		return room.RoomInfo{}, fmt.Errorf("json parsing err : %w", err)
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

func (r *roomRepository) Update(ctx context.Context, key string, room *room.RoomInfo) error {

	if err := r.rdb.HSet(ctx, key, room.ConvertRedisData()); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	if err := r.rdb.Expire(ctx, key, time.Duration(room.RoomIdTTLDay*24)*time.Hour); err != nil {
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

func (r *roomRepository) SetRoomMap(ctx context.Context, key string, data *room.RoomInfo) error {

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

func (r *roomRepository) GetRoomMap(ctx context.Context, key, mapKey string) (room.RoomInfo, error) {

	result, err := r.rdb.HGet(ctx, key, mapKey)
	if err != nil {
		return room.RoomInfo{}, err
	}

	var roomInfo room.RoomInfo
	if err := json.Unmarshal([]byte(result), &roomInfo); err != nil {
		return room.RoomInfo{}, fmt.Errorf("json parsing err : %w", err)
	}

	return roomInfo, nil
}
