package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"live-chat-server/internal/database/mysql"
	"live-chat-server/internal/domain/system"
	"strings"
)

// repository 에서는 쿼리를 작성하고, struct 으로 변경한다.
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

	rows, err := s.db.GetServerList(qs)
	if err != nil {
		return nil, err
	}

	var result []system.ServerInfo

	for _, row := range rows {
		serverInfo := system.ServerInfo{}

		if ip, ok := row["ip"].(string); ok {
			serverInfo.IP = ip
		}

		if available, ok := row["available"].(int); ok {
			isAvailable := available == 1
			serverInfo.Available = isAvailable
		}

		result = append(result, serverInfo)
	}

	return result, nil
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}
