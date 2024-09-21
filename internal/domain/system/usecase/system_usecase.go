package usecase

import (
	"live-chat-server/internal/domain/system"
	"log"
)

type systemUseCase struct {
	systemRepo    system.Repository
	systemPubSub  system.PubSub
	avgServerList map[string]bool
}

func NewSystemUseCase(repository system.Repository, systemPubSub system.PubSub) system.UseCase {

	s := &systemUseCase{
		systemRepo:    repository,
		systemPubSub:  systemPubSub,
		avgServerList: make(map[string]bool),
	}

	if err := s.setServerInfo(); err != nil {
		log.Fatalf("failed register server info, err : %v", err)
	}

	if err := s.systemPubSub.RegisterSubTopic("chat"); err != nil {
		log.Fatalf("failed register topic, err : %v", err)
	}

	return s
}

func (s *systemUseCase) GetServerList() ([]system.ServerInfo, error) {

	if len(s.avgServerList) == 0 {
		return []system.ServerInfo{}, nil
	}

	var res []system.ServerInfo

	for ip, available := range s.avgServerList {
		if available {
			server := system.ServerInfo{
				IP: ip,
			}
			res = append(res, server)
		}
	}

	return res, nil
}

func (s *systemUseCase) setServerInfo() error {

	serverList, err := s.GetAvailableServerList()
	if err != nil {
		return err
	}

	for _, server := range serverList {
		s.avgServerList[server.IP] = true
	}

	return nil
}

func (s *systemUseCase) GetAvailableServerList() ([]system.ServerInfo, error) {
	return s.systemRepo.GetAvailableServerList()
}
