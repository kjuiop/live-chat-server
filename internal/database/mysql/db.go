package mysql

import "live-chat-server/internal/domain/system"

type Client interface {
	GetServerList(qs string) ([]system.ServerInfo, error)
}
