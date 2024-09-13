package mocks

import "live-chat-server/internal/domain/system"

type SystemUseCaseStub struct {
}

func NewSystemUseCaseStub() system.UseCase {
	return &SystemUseCaseStub{}
}
