package kafka

import "live-chat-server/internal/mq/types"

type memoryClient struct {
}

func NewMemoryClient() Client {
	return &memoryClient{}
}

func (m memoryClient) Subscribe(topic string) error {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) Poll(timeoutMs int) types.Event {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) PublishEvent(topic string, data []byte) (types.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (m memoryClient) Close(mqType string) {
	//TODO implement me
	panic("implement me")
}
