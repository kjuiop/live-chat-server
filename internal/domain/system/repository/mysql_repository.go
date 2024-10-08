package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"live-chat-server/internal/database/mysql"
	"live-chat-server/internal/domain/system"
	"log/slog"
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

func (s *systemMySqlRepository) SetChatServerInfo(ip string, available bool) error {
	qs := query([]string{
		"INSERT INTO",
		"chatting.serverInfo(`ip`, `available`)",
		"VALUES (?, ?)",
		"ON DUPLICATE KEY UPDATE `available` = VALUES(`available`)",
	})
	return s.db.ExecQuery(qs, ip, available)
}

func (s *systemMySqlRepository) GetAvailableServerList() ([]system.ServerInfo, error) {
	qs := query([]string{"SELECT * FROM", serverInfo, "WHERE available = 1"})

	rows, err := s.db.ExecQueryAndFetchRows(qs)
	if err != nil {
		return nil, err
	}

	var result []system.ServerInfo

	for _, row := range rows {
		serverInfo := system.ServerInfo{}

		if ipBytes, ok := row["ip"].([]byte); ok {
			serverInfo.IP = string(ipBytes) // []byte를 string으로 변환
		} else if ip, ok := row["ip"].(string); ok {
			serverInfo.IP = ip
		} else {
			slog.Error("unexpected type for ip field", "value", row["ip"])
			continue // 타입이 예상과 다르면 현재 반복을 종료하고 다음으로 넘어감
		}

		if available, ok := row["available"].(int64); ok {
			isAvailable := available == 1
			serverInfo.Available = isAvailable
		} else {
			slog.Error("unexpected type for available field", "value", row["available"])
			continue
		}

		result = append(result, serverInfo)
	}

	return result, nil
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}
