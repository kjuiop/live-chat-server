package system

type HealthRes struct {
	Message string `json:"message"`
}

type ServerInfo struct {
	IP        string `json:"ip"`
	Available bool   `json:"available"`
}

type UseCase interface {
	GetServerList() ([]ServerInfo, error)
}

type Repository interface {
	GetAvailableServerList() ([]ServerInfo, error)
}

type PubSub interface {
}
