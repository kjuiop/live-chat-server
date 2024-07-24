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
