package mocks

import (
	"context"
	"github.com/gorilla/websocket"
	"live-chat-server/internal/domain/chat"
	"live-chat-server/internal/domain/room"
	"sync"
)

type ChatUseCaseStub struct {
	crMutex *sync.RWMutex
	hub     map[string]*chat.Room
}

func NewChatUseCaseStub() chat.ChatUseCase {
	return &ChatUseCaseStub{
		crMutex: &sync.RWMutex{},
		hub:     make(map[string]*chat.Room),
	}
}

func (cs *ChatUseCaseStub) ServeWs(ctx context.Context, socket *websocket.Conn, chatRoom *chat.Room, userId string) error {
	client := chat.NewClient(socket, chatRoom, userId)

	chatRoom.Join <- client

	defer func() {
		chatRoom.Leave <- client
	}()

	go client.Write()

	client.Read()

	return nil
}

func (cs *ChatUseCaseStub) GetChatRoom(ctx context.Context, roomId string) (*chat.Room, error) {
	cs.crMutex.Lock()
	defer cs.crMutex.Unlock()

	if _, ok := cs.hub[roomId]; !ok {
		roomInfo := room.RoomInfo{RoomId: roomId}
		cs.hub[roomId] = chat.NewChatRoom(roomInfo)
	}

	return cs.hub[roomId], nil
}
