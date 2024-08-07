package models

type ChatRoom struct {
	RoomId string `json:"RoomId"`
}

func NewChatRoom(roomInfo RoomInfo) *ChatRoom {
	return &ChatRoom{
		RoomId: roomInfo.RoomId,
	}
}
