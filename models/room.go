package models

import (
	"fmt"
	"live-chat-server/utils"
	"math/rand"
	"strings"
	"time"
)

type CreateRoomReq struct {
	CustomerId   string `json:"customer_id"`
	ChannelKey   string `json:"channel_key"`
	BroadCastKey string `json:"broadcast_key"`
}

type RoomInfo struct {
	RoomId       string `json:"room_id"`
	CustomerId   string `json:"customer_id"`
	ChannelKey   string `json:"channel_key"`
	BroadCastKey string `json:"broadcast_key"`
}

func NewRoomInfo(req *CreateRoomReq, prefix string) *RoomInfo {
	return &RoomInfo{
		RoomId:       fmt.Sprintf("%s_%s", getChatPrefix(prefix), utils.GenUUID()),
		CustomerId:   req.CustomerId,
		ChannelKey:   req.ChannelKey,
		BroadCastKey: req.BroadCastKey,
	}
}

func getChatPrefix(prefix string) string {
	array := strings.Split(prefix, ",")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[rand.Intn(len(array))]
}
