package pubsub

import (
	"live-chat-server/config"
	"live-chat-server/internal/domain/system"
	"live-chat-server/internal/mq/kafka"
)

type PubSub struct {
	cfg config.Kafka
	mq  kafka.Client
}

func NewSystemPubSub(cfg config.Kafka, mq kafka.Client) system.PubSub {
	return &PubSub{
		cfg: cfg,
		mq:  mq,
	}
}
