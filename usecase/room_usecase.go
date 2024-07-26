package usecase

import (
	"context"
	"live-chat-server/models"
	"time"
)

type roomUseCase struct {
	roomRepository models.RoomRepository
}

func NewRoomUseCase(roomRepository models.RoomRepository) models.RoomUseCase {
	return &roomUseCase{
		roomRepository: roomRepository,
	}
}

func (r *roomUseCase) CreateChatRoom(c context.Context, room *models.RoomInfo) error {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	return r.roomRepository.Create(ctx, room)
}

func (r *roomUseCase) GetChatRoomById(c context.Context, roomId string) (models.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	roomInfo, err := r.roomRepository.Fetch(ctx, roomId)
	if err != nil {
		return models.RoomInfo{}, err
	}

	return roomInfo, nil
}

func (r *roomUseCase) CheckExistRoomId(c context.Context, roomId string) (bool, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	isExist, err := r.roomRepository.Exists(ctx, roomId)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (r *roomUseCase) UpdateChatRoom(c context.Context, roomId string, room *models.RoomInfo) (models.RoomInfo, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
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
