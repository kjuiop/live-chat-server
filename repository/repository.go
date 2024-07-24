package repository

import (
	"context"
	"live-chat-server/models"
)

type Client interface {
	CreateChatRoom(ctx context.Context, room *models.RoomInfo) error
}
