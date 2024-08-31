package usecase

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"live-chat-server/domain/chat"
	"live-chat-server/domain/chat/types"
	"live-chat-server/domain/room"
	"net/http"
	"sync"
	"time"
)

var crMutex = &sync.RWMutex{}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  types.SocketBufferSize,
	WriteBufferSize: types.MessageBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type chatUseCase struct {
	roomUseCase    room.RoomUseCase
	contextTimeout time.Duration
	hub            map[string]*chat.Room
}

func NewChatUseCase(roomUseCase room.RoomUseCase, timeout time.Duration) chat.ChatUseCase {
	return &chatUseCase{
		roomUseCase:    roomUseCase,
		contextTimeout: timeout,
		hub:            make(map[string]*chat.Room),
	}
}

func (cu *chatUseCase) ServeWs(c context.Context, writer http.ResponseWriter, request *http.Request, roomId, userId string) error {

	chatRoom, err := cu.getChatRoom(c, roomId)
	if err != nil {
		return err
	}

	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return fmt.Errorf("failed connect socket, err : %w", err)
	}

	client := chat.NewClient(socket, chatRoom, userId)

	chatRoom.Join <- client

	defer func() {
		chatRoom.Leave <- client
	}()

	go client.Write()

	client.Read()

	return nil
}

func (cu *chatUseCase) getChatRoom(c context.Context, roomId string) (*chat.Room, error) {

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
