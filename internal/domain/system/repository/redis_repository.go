package repository

import (
	"live-chat-server/config"
	"live-chat-server/internal/domain/system"
)

type systemRedisRepository struct {
	cfg config.Redis
}

func NewSystemRedisRepository(cfg config.Redis) system.Repository {
	return nil
}
