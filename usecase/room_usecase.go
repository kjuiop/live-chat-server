package usecase

import (
	"context"
	"live-chat-server/models"
	"time"
)

type roomUseCase struct {
	roomRepository models.RoomRepository
	contextTimeout time.Duration
}

func NewRoomUseCase(roomRepository models.RoomRepository, timeout time.Duration) models.RoomUseCase {
	return &roomUseCase{
		roomRepository: roomRepository,
		contextTimeout: timeout,
	}
}

func (r *roomUseCase) CreateChatRoom(c context.Context, room *models.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	return r.roomRepository.Create(ctx, room)
}

func (r *roomUseCase) GetChatRoomById(c context.Context, roomId string) (models.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	roomInfo, err := r.roomRepository.Fetch(ctx, roomId)
	if err != nil {
		return models.RoomInfo{}, err
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

func (r *roomUseCase) UpdateChatRoom(c context.Context, roomId string, room *models.RoomInfo) (models.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	if err := r.roomRepository.Update(ctx, roomId, room); err != nil {
		return models.RoomInfo{}, err
	}

	savedInfo, err := r.roomRepository.Fetch(c, roomId)
	if err != nil {
		return models.RoomInfo{}, err
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
