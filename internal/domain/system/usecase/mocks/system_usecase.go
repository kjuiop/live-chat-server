package mocks

import "live-chat-server/internal/domain/system"

type SystemUseCaseStub struct {
}

func NewSystemUseCaseStub() system.UseCase {
	return &SystemUseCaseStub{}
}

func (s SystemUseCaseStub) GetServerList() ([]system.ServerInfo, error) {
	//TODO implement me
	panic("implement me")
}
