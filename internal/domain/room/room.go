package room

import (
	"context"
	"fmt"
	"live-chat-server/internal/utils"
	"math/rand"
	"strings"
	"time"
)

type RoomRequest struct {
	CustomerId   string `json:"customer_id" binding:"required"`
	ChannelKey   string `json:"channel_key" binding:"required"`
	BroadCastKey string `json:"broadcast_key" binding:"required"`
}

type RoomIdRequest struct {
	ChannelKey   string `form:"channel_key" binding:"required"`
	BroadCastKey string `form:"broadcast_key" binding:"required"`
}

type RoomResponse struct {
	RoomId       string `json:"room_id"`
	CustomerId   string `json:"customer_id,omitempty"`
	ChannelKey   string `json:"channel_key,omitempty"`
	BroadcastKey string `json:"broadcast_key,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
}

type RoomInfo struct {
	RoomId       string `json:"room_id"`
	CustomerId   string `json:"customer_id"`
	ChannelKey   string `json:"channel_key"`
	BroadcastKey string `json:"broadcast_key"`
	CreatedAt    int64  `json:"created_at"`
}

func NewRoomInfo(req RoomRequest, prefix string) *RoomInfo {
	return &RoomInfo{
		RoomId:       fmt.Sprintf("%s-%s", getChatPrefix(prefix), utils.GenUUID()),
		CustomerId:   req.CustomerId,
		ChannelKey:   req.ChannelKey,
		BroadcastKey: req.BroadCastKey,
		CreatedAt:    time.Now().Unix(),
	}
}

func UpdateRoomInfo(req RoomRequest, roomId string) *RoomInfo {
	return &RoomInfo{
		RoomId:       roomId,
		CustomerId:   req.CustomerId,
		ChannelKey:   req.ChannelKey,
		BroadcastKey: req.BroadCastKey,
		CreatedAt:    time.Now().Unix(),
	}
}

func (r *RoomInfo) ConvertRedisData() map[string]interface{} {
	return map[string]interface{}{
		"room_id":       r.RoomId,
		"customer_id":   r.CustomerId,
		"channel_key":   r.ChannelKey,
		"broadcast_key": r.BroadcastKey,
		"created_at":    r.CreatedAt,
	}
}

func getChatPrefix(prefix string) string {
	array := strings.Split(prefix, ",")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[rand.Intn(len(array))]
}

type RoomUseCase interface {
	CreateChatRoom(ctx context.Context, room RoomInfo) error
	GetChatRoomById(ctx context.Context, roomId string) (RoomInfo, error)
	CheckExistRoomId(ctx context.Context, roomId string) (bool, error)
	UpdateChatRoom(ctx context.Context, roomId string, room RoomInfo) (RoomInfo, error)
	DeleteChatRoom(ctx context.Context, roomId string) error
	RegisterRoomId(ctx context.Context, room RoomInfo) error
	GetChatRoomId(ctx context.Context, room RoomIdRequest) (RoomInfo, error)
}

type RoomRepository interface {
	Create(ctx context.Context, data RoomInfo) error
	Fetch(ctx context.Context, roomId string) (RoomInfo, error)
	Exists(ctx context.Context, roomId string) (bool, error)
	Update(ctx context.Context, roomId string, data RoomInfo) error
	Delete(ctx context.Context, roomId string) error
	SetRoomMap(ctx context.Context, data RoomInfo) error
	GetRoomMap(ctx context.Context, channelKey, broadcastKey string) (RoomInfo, error)
}
