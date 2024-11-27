package kafka

import "live-chat-server/internal/mq/types"

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}

func (m memoryClient) Subscribe(topic string) error {
	return nil
}

func (m memoryClient) Poll(timeoutMs int) types.Event {
	return nil
}

func (m memoryClient) PublishEvent(topic string, data []byte) (types.Event, error) {
	return nil, nil
}
