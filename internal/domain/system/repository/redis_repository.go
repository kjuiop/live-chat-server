package repository

import (
	"context"
	"fmt"
	"live-chat-server/internal/database/redis"
	"live-chat-server/internal/domain/system"
)

const (
	LiveChatServerInfo = "live-chat-server-info"
)

type systemRedisRepository struct {
	db redis.Client
}

func NewSystemRedisRepository(db redis.Client) system.Repository {
	return &systemRedisRepository{
		db: db,
	}
}

func (s *systemRedisRepository) SetChatServerInfo(ip string, available bool) error {

	data := system.NewServerInfo(ip, available)
	if err := s.db.HSet(context.TODO(), LiveChatServerInfo, data.IP, data.ConvertRedisData()); err != nil {
		return err
	}

	return nil
}

func (s *systemRedisRepository) GetAvailableServerList() ([]system.ServerInfo, error) {
	result, err := s.db.HGetAll(context.TODO(), LiveChatServerInfo)
	if err != nil {
		return nil, fmt.Errorf("fail get room info, err : %w", err)
	}

	// redis 는 하나의 json 문자열로 반환하기 때문에 안의 데이터타입을 변경하기 위해서는 아래와 같은 매핑 작업을 필요로 함
	data := make(map[string]interface{})
	for k, v := range result {
		data[k] = v
	}

	return nil, nil
}
