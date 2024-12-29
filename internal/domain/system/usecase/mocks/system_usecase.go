package mocks

import (
	"live-chat-server/internal/domain/system"
	"live-chat-server/internal/mq/types"
)

type SystemUseCaseStub struct {
}

func (s SystemUseCaseStub) GetServerList() ([]system.ServerInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemUseCaseStub) RegisterSubTopic(topic string) error {
	//TODO implement me
	panic("implement me")
}

func (s SystemUseCaseStub) SetChatServerInfo(ip string, available bool) error {
	//TODO implement me
	panic("implement me")
}

func (s SystemUseCaseStub) PublishServerStatusEvent(addr string, status bool) {
	//TODO implement me
	panic("implement me")
}

func (s SystemUseCaseStub) LoopSubKafka(timeoutMs int) (*types.Message, error) {
	//TODO implement me
	panic("implement me")
}

func NewSystemUseCaseStub() system.UseCase {
	return &SystemUseCaseStub{}
}
