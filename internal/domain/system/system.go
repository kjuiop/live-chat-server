package system

import (
	"live-chat-server/internal/mq/types"
)

type HealthRes struct {
	Message string `json:"message"`
}

type ServerInfo struct {
	IP        string `json:"ip"`
	Available bool   `json:"available"`
}

type UseCase interface {
	GetServerList() ([]ServerInfo, error)
	SetChatServerInfo(ip string, available bool) error
	PublishServerStatusEvent(addr string, status bool)
}

type Repository interface {
	GetAvailableServerList() ([]ServerInfo, error)
	SetChatServerInfo(ip string, available bool) error
}

type PubSub interface {
	RegisterSubTopic(topic string) error
	Poll(timeoutMs int) types.Event
	PublishEvent(topic string, data []byte) (types.Event, error)
}
