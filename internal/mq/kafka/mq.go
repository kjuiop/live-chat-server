package kafka

import "live-chat-server/internal/mq/types"

type Client interface {
	Subscribe(topic string) error
	Poll(timeoutMs int) types.Event
}
