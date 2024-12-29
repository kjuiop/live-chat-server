package usecase

import (
	"context"
	"live-chat-server/internal/domain/room"
	"time"
)

type roomUseCase struct {
	roomRepository room.RoomRepository
	contextTimeout time.Duration
}

func NewRoomUseCase(roomRepository room.RoomRepository, timeout time.Duration) room.RoomUseCase {
	return &roomUseCase{
		roomRepository: roomRepository,
		contextTimeout: timeout,
	}
}

func (r *roomUseCase) CreateChatRoom(c context.Context, room room.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	return r.roomRepository.Create(ctx, room)
}

func (r *roomUseCase) GetChatRoomById(c context.Context, roomId string) (room.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	roomInfo, err := r.roomRepository.Fetch(ctx, roomId)
	if err != nil {
		return room.RoomInfo{}, err
	}

	return roomInfo, nil
}

func (r *roomUseCase) CheckExistRoomId(c context.Context, roomId string) (bool, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	isExist, err := r.roomRepository.Exists(ctx, roomId)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (r *roomUseCase) UpdateChatRoom(c context.Context, roomId string, roomInfo room.RoomInfo) (room.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepository.Update(ctx, roomId, roomInfo); err != nil {
		return room.RoomInfo{}, err
	}

	savedInfo, err := r.roomRepository.Fetch(c, roomId)
	if err != nil {
		return room.RoomInfo{}, err
	}

	return savedInfo, nil
}

func (r *roomUseCase) DeleteChatRoom(c context.Context, roomId string) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepository.Delete(ctx, roomId); err != nil {
		return err
	}

	return nil
}

func (r *roomUseCase) RegisterRoomId(c context.Context, roomInfo room.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepository.SetRoomMap(ctx, roomInfo); err != nil {
		return err
	}

	return nil
}

func (r *roomUseCase) GetChatRoomId(c context.Context, req room.RoomIdRequest) (room.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	roomInfo, err := r.roomRepository.GetRoomMap(ctx, req.ChannelKey, req.BroadCastKey)
	if err != nil {
		return room.RoomInfo{}, err
	}

	return roomInfo, nil
}
