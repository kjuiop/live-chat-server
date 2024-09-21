package pubsub

import (
	"live-chat-server/config"
	"live-chat-server/internal/domain/system"
	"live-chat-server/internal/mq/kafka"
	"live-chat-server/internal/mq/types"
)

type MemoryPubSub struct {
	cfg config.Kafka
	mq  kafka.Client
}

func NewMemorySystemPubSub(cfg config.Kafka, mq kafka.Client) system.PubSub {
	return &PubSub{
		cfg: cfg,
		mq:  mq,
	}
}

func (m *MemoryPubSub) RegisterSubTopic(topic string) error {
	if err := m.mq.Subscribe(topic); err != nil {
		return err
	}
	return nil
}

func (m *MemoryPubSub) Poll(timeoutMs int) types.Event {
	return m.mq.Poll(timeoutMs)
}
