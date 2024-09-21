package usecase

import (
	"encoding/json"
	"live-chat-server/internal/domain/system"
	"live-chat-server/internal/mq/types"
	"log"
	"log/slog"
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

	go s.loopSubKafka()

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

func (s *systemUseCase) loopSubKafka() {
	for {
		ev := s.systemPubSub.Poll(100)
		if ev.IsError() {
			errorEvent := ev.(*types.Error)
			slog.Error("Failed to Polling event", "error", errorEvent.Error)
			continue
		}

		if ev.IsMessage() {
			message := ev.(*types.Message)

			var decoder system.ServerInfo
			if err := json.Unmarshal(message.Value, &decoder); err != nil {
				slog.Error("failed to decode event", "event_value", string(message.Value))
				continue
			}

			slog.Debug("received kafka event", "event_value", string(message.Value))
			s.avgServerList[decoder.IP] = decoder.Available
		}
	}
}

func (s *systemUseCase) GetAvailableServerList() ([]system.ServerInfo, error) {
	return s.systemRepo.GetAvailableServerList()
}
