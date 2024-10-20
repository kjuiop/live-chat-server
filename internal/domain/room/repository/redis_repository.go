package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"live-chat-server/internal/database/redis"
	"live-chat-server/internal/domain/room"
	"time"
)

const (
	LiveChatServerRoom    = "live-chat-server-room"
	LiveChatServerRoomMap = "live-chat-server-room-map"
)

const (
	RoomExpire = time.Duration(7) * 24 * time.Hour
)

// 실제 쿼리가 작성되는 공간

type roomRedisRepository struct {
	db redis.Client
}

func NewRoomRedisRepository(client redis.Client) room.RoomRepository {
	return &roomRedisRepository{
		db: client,
	}
}

func (r *roomRedisRepository) Create(ctx context.Context, room room.RoomInfo) error {

	if err := r.db.Set(ctx, convertRoomKey(room.RoomId), room.ConvertRedisData(), RoomExpire); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	return nil
}

func (r *roomRedisRepository) Fetch(ctx context.Context, roomId string) (room.RoomInfo, error) {

	result, err := r.db.Get(ctx, convertRoomKey(roomId))
	if err != nil {
		return room.RoomInfo{}, fmt.Errorf("fail get room info, err : %w", err)
	}

	// redis 는 하나의 json 문자열로 반환하기 때문에 안의 데이터타입을 변경하기 위해서는 아래와 같은 매핑 작업을 필요로 함
	//data := make(map[string]interface{})
	//for k, v := range result {
	//	if k == "created_at" {
	//		createdAt, _ := strconv.Atoi(v)
	//		data[k] = createdAt
	//		continue
	//	} else {
	//		data[k] = v
	//	}
	//}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return room.RoomInfo{}, fmt.Errorf("json marshal err : %w", err)
	}

	var roomInfo room.RoomInfo
	if err := json.Unmarshal(jsonData, &roomInfo); err != nil {
		return room.RoomInfo{}, fmt.Errorf("json parsing err : %w", err)
	}

	return roomInfo, nil
}

func (r *roomRedisRepository) Exists(ctx context.Context, roomId string) (bool, error) {

	isExist, err := r.db.Exists(ctx, convertRoomKey(roomId))
	if err != nil {
		return false, fmt.Errorf("fail redis cmd exist err : %w", err)
	}

	return isExist, nil
}

func (r *roomRedisRepository) Update(ctx context.Context, roomId string, room room.RoomInfo) error {

	if err := r.db.Set(ctx, convertRoomKey(roomId), room.ConvertRedisData(), RoomExpire); err != nil {
		return fmt.Errorf("create chat room hm set err : %w", err)
	}

	return nil
}

func (r *roomRedisRepository) Delete(ctx context.Context, roomId string) error {

	if err := r.db.DelByKey(ctx, convertRoomKey(roomId)); err != nil {
		return err
	}

	return nil
}

func (r *roomRedisRepository) SetRoomMap(ctx context.Context, data room.RoomInfo) error {

	jData, err := json.Marshal(data.ConvertRedisData())
	if err != nil {
		return fmt.Errorf("set room map json encoding fail, err : %w", err)
	}

	if err := r.db.Set(ctx, generateRoomMapKey(data.ChannelKey, data.BroadcastKey), string(jData), RoomExpire); err != nil {
		return fmt.Errorf("set room map err : %w", err)
	}

	return nil
}

func (r *roomRedisRepository) GetRoomMap(ctx context.Context, key, mapKey string) (room.RoomInfo, error) {

	result, err := r.db.HGet(ctx, key, mapKey)
	if err != nil {
		return room.RoomInfo{}, err
	}

	var roomInfo room.RoomInfo
	if err := json.Unmarshal([]byte(result), &roomInfo); err != nil {
		return room.RoomInfo{}, fmt.Errorf("json parsing err : %w", err)
	}

	return roomInfo, nil
}

func convertRoomKey(roomId string) string {
	return fmt.Sprintf("%s_%s", LiveChatServerRoom, roomId)
}

func generateRoomMapKey(channelKey, broadcastKey string) string {
	return fmt.Sprintf("%s_%s_%s", LiveChatServerRoomMap, channelKey, broadcastKey)
}
