package repository

import (
	"context"
	"live-chat-server/models"
)

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}

func (m memoryClient) CreateChatRoom(ctx context.Context, room *models.RoomInfo) error {
	return nil
}
