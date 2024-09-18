package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"live-chat-server/internal/database/mysql"
	"live-chat-server/internal/domain/system"
	"strings"
)

type systemMySqlRepository struct {
	db mysql.Client
}

const (
	serverInfo = "chatting.serverInfo"
)

func NewSystemMySqlRepository(db mysql.Client) system.Repository {
	return &systemMySqlRepository{
		db: db,
	}
}

func (s *systemMySqlRepository) GetAvailableServerList() ([]system.ServerInfo, error) {
	qs := query([]string{"SELECT * FROM", serverInfo, "WHERE available = 1"})

	list, err := s.db.GetServerList(qs)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}
