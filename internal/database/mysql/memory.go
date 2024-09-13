package mysql

import "live-chat-server/internal/domain/system"

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}

func (m memoryClient) GetServerList(qs string) ([]system.ServerInfo, error) {
	//TODO implement me
	panic("implement me")
}
