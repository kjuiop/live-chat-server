package models

import (
	"context"
	"fmt"
	"live-chat-server/utils"
	"math/rand"
	"strings"
	"time"
)

type RoomRequest struct {
	CustomerId   string `json:"customer_id"`
	ChannelKey   string `json:"channel_key"`
	BroadCastKey string `json:"broadcast_key"`
}

type RoomInfo struct {
	RoomId       string `json:"room_id"`
	CustomerId   string `json:"customer_id"`
	ChannelKey   string `json:"channel_key"`
	BroadcastKey string `json:"broadcast_key"`
}

func NewRoomInfo(req *RoomRequest, prefix string) *RoomInfo {
	return &RoomInfo{
		RoomId:       fmt.Sprintf("%s_%s", getChatPrefix(prefix), utils.GenUUID()),
		CustomerId:   req.CustomerId,
		ChannelKey:   req.ChannelKey,
		BroadcastKey: req.BroadCastKey,
	}
}

func UpdateRoomInfo(req *RoomRequest, roomId string) *RoomInfo {
	return &RoomInfo{
		RoomId:       roomId,
		CustomerId:   req.CustomerId,
		ChannelKey:   req.ChannelKey,
		BroadcastKey: req.BroadCastKey,
	}
}

func (r *RoomInfo) ConvertRedisData() map[string]interface{} {
	return map[string]interface{}{
		"room_id":       r.RoomId,
		"customer_id":   r.CustomerId,
		"channel_key":   r.ChannelKey,
		"broadcast_key": r.BroadcastKey,
	}
}

func getChatPrefix(prefix string) string {
	array := strings.Split(prefix, ",")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[rand.Intn(len(array))]
}

type RoomUseCase interface {
	CreateChatRoom(ctx context.Context, room *RoomInfo) error
	GetChatRoomById(ctx context.Context, roomId string) (RoomInfo, error)
	CheckExistRoomId(ctx context.Context, roomId string) (bool, error)
	UpdateChatRoom(ctx context.Context, roomId string, room *RoomInfo) (RoomInfo, error)
}

type RoomRepository interface {
	Create(ctx context.Context, data *RoomInfo) error
	Fetch(ctx context.Context, key string) (RoomInfo, error)
	Exists(ctx context.Context, key string) (bool, error)
	Update(ctx context.Context, key string, data *RoomInfo) error
}
