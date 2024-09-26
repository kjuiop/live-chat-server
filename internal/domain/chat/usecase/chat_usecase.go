package usecase

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"live-chat-server/internal/domain/chat"
	"live-chat-server/internal/domain/chat/types"
	"live-chat-server/internal/domain/room"
	"net/http"
	"sync"
	"time"
)

var crMutex = &sync.RWMutex{}

type chatUseCase struct {
	roomUseCase    room.RoomUseCase
	contextTimeout time.Duration
	hub            map[string]*chat.Room
	upgrader       *websocket.Upgrader
}

func NewChatUseCase(roomUseCase room.RoomUseCase, timeout time.Duration) chat.ChatUseCase {
	return &chatUseCase{
		roomUseCase:    roomUseCase,
		contextTimeout: timeout,
		hub:            make(map[string]*chat.Room),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  types.SocketBufferSize,
			WriteBufferSize: types.MessageBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (cu *chatUseCase) ServeWs(c context.Context, socket *websocket.Conn, chatRoom *chat.Room, userId string) error {

	client := chat.NewClient(socket, chatRoom, userId)

	chatRoom.Join <- client

	defer func() {
		chatRoom.Leave <- client
	}()

	go client.Write()

	client.Read()

	return nil
}

func (cu *chatUseCase) GetChatRoom(c context.Context, roomId string) (*chat.Room, error) {

	crMutex.Lock()
	defer func() {
		crMutex.Unlock()
	}()

	if _, ok := cu.hub[roomId]; !ok {
		roomInfo, err := cu.roomUseCase.GetChatRoomById(c, roomId)
		if err != nil {
			return nil, fmt.Errorf("not found chat room, key : %s, err : %w", roomId, err)
		}
		cu.hub[roomId] = chat.NewChatRoom(roomInfo)
	}

	return cu.hub[roomId], nil
}
