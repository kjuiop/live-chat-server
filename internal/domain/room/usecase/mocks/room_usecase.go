package mocks

import (
	"context"
	"errors"
	"fmt"
	"live-chat-server/internal/domain/room"
	"sync"
)

type RoomUseCaseStub struct {
	mu    sync.Mutex
	Rooms map[string]room.RoomInfo
}

func NewRoomUseCaseStub(rm []room.RoomInfo) room.RoomUseCase {
	stub := &RoomUseCaseStub{
		Rooms: make(map[string]room.RoomInfo),
	}

	for _, v := range rm {
		stub.Rooms[v.RoomId] = v
	}

	return stub
}

func (s *RoomUseCaseStub) CreateChatRoom(ctx context.Context, room room.RoomInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Rooms[room.RoomId]; exists {
		return errors.New("room already exists")
	}

	s.Rooms[room.RoomId] = room
	return nil
}

func (s *RoomUseCaseStub) GetChatRoomById(ctx context.Context, roomId string) (room.RoomInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomInfo, exists := s.Rooms[roomId]
	if !exists {
		return room.RoomInfo{}, fmt.Errorf("room not found, room_id : %s", roomId)
	}

	return roomInfo, nil
}

func (s *RoomUseCaseStub) CheckExistRoomId(ctx context.Context, roomId string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.Rooms[roomId]
	return exists, nil
}

func (s *RoomUseCaseStub) UpdateChatRoom(ctx context.Context, roomId string, roomInfo room.RoomInfo) (room.RoomInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Rooms[roomId]; !exists {
		return room.RoomInfo{}, errors.New("room not found")
	}

	s.Rooms[roomId] = roomInfo
	return roomInfo, nil
}

func (s *RoomUseCaseStub) DeleteChatRoom(ctx context.Context, roomId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Rooms[roomId]; !exists {
		return errors.New("room not found")
	}

	delete(s.Rooms, roomId)
	return nil
}

func (s *RoomUseCaseStub) RegisterRoomId(ctx context.Context, roomInfo room.RoomInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	//if _, exists := s.Rooms[roomInfo.RoomId]; exists {
	//	return errors.New("room already registered")
	//}

	s.Rooms[roomInfo.RoomId] = roomInfo
	return nil
}

func (s *RoomUseCaseStub) GetChatRoomId(ctx context.Context, req room.RoomIdRequest) (room.RoomInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, roomInfo := range s.Rooms {
		if roomInfo.ChannelKey == req.ChannelKey && roomInfo.BroadcastKey == req.BroadCastKey {
			return roomInfo, nil
		}
	}

	return room.RoomInfo{}, errors.New("room not found")
}
