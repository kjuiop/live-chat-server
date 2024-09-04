package mysql

import (
	"context"
	"live-chat-server/config"
)

type mysqlClient struct {
	cfg config.Mysql
}

func NewMysqlSingleClient(ctx context.Context, cfg config.Mysql) (Client, error) {
	return &mysqlClient{
		cfg: cfg,
	}, nil
}
