package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"live-chat-server/config"
	"live-chat-server/internal/domain/system"
)

type mysqlClient struct {
	cfg config.Mysql
	db  *sql.DB
}

func NewMysqlSingleClient(ctx context.Context, cfg config.Mysql) (Client, error) {

	db, err := sql.Open(cfg.Database, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed connect mysql, err : %w", err)
	}

	return &mysqlClient{
		cfg: cfg,
		db:  db,
	}, nil
}

func (m *mysqlClient) GetServerList(qs string) ([]system.ServerInfo, error) {

	if cursor, err := m.db.Query(qs); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []system.ServerInfo

		for cursor.Next() {
			d := new(system.ServerInfo)

			if err = cursor.Scan(
				d.IP,
				d.Available,
			); err != nil {
				return nil, err
			} else {
				result = append(result, *d)
			}
		}

		if len(result) == 0 {
			return []system.ServerInfo{}, nil
		} else {
			return result, nil
		}
	}
}
