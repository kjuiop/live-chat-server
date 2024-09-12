package repository

import (
	"live-chat-server/config"
	"live-chat-server/internal/domain/system"
)

type systemMySqlRepository struct {
	cfg config.Mysql
}

func NewSystemMySqlRepository(cfg config.Mysql) system.Repository {
	return nil
}
