package usecase

import (
	"context"
	"fmt"
	"live-chat-server/domain/room"
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

	roomInfo, err := r.roomRepository.Fetch(ctx, getConvertRedisKey(roomId))
	if err != nil {
		return room.RoomInfo{}, err
	}

	return roomInfo, nil
}

func (r *roomUseCase) CheckExistRoomId(c context.Context, roomId string) (bool, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	isExist, err := r.roomRepository.Exists(ctx, getConvertRedisKey(roomId))
	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (r *roomUseCase) UpdateChatRoom(c context.Context, roomId string, roomInfo room.RoomInfo) (room.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepository.Update(ctx, getConvertRedisKey(roomId), roomInfo); err != nil {
		return room.RoomInfo{}, err
	}

	savedInfo, err := r.roomRepository.Fetch(c, getConvertRedisKey(roomId))
	if err != nil {
		return room.RoomInfo{}, err
	}

	return savedInfo, nil
}

func (r *roomUseCase) DeleteChatRoom(c context.Context, roomId string) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepository.Delete(ctx, getConvertRedisKey(roomId)); err != nil {
		return err
	}

	return nil
}

func (r *roomUseCase) RegisterRoomId(c context.Context, roomInfo room.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepository.SetRoomMap(ctx, room.LiveChatServerRoomList, roomInfo); err != nil {
		return fmt.Errorf("failed create room map, key:%s, err : %w", room.LiveChatServerRoomList, err)
	}

	return nil
}

func (r *roomUseCase) GetChatRoomId(c context.Context, req room.RoomIdRequest) (room.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	roomMapKey := fmt.Sprintf("%s_%s_%s", room.LiveChatServerRoomList, req.ChannelKey, req.BroadCastKey)

	roomInfo, err := r.roomRepository.GetRoomMap(ctx, room.LiveChatServerRoomList, roomMapKey)
	if err != nil {
		return room.RoomInfo{}, err
	}

	return roomInfo, nil
}

func getConvertRedisKey(roomId string) string {
	return fmt.Sprintf("%s_%s", room.LiveChatServerRoom, roomId)
}
