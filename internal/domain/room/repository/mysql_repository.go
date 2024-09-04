package repository

import (
	"context"
	"live-chat-server/internal/database/mysql"
	"live-chat-server/internal/domain/room"
)

type roomMysqlRepository struct {
	db mysql.Client
}

func NewRoomMysqlRepository(client mysql.Client) room.RoomRepository {
	return &roomMysqlRepository{
		db: client,
	}
}

func (r roomMysqlRepository) Create(ctx context.Context, data room.RoomInfo) error {
	//TODO implement me
	panic("implement me")
}

func (r roomMysqlRepository) Fetch(ctx context.Context, key string) (room.RoomInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r roomMysqlRepository) Exists(ctx context.Context, key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r roomMysqlRepository) Update(ctx context.Context, key string, data room.RoomInfo) error {
	//TODO implement me
	panic("implement me")
}

func (r roomMysqlRepository) Delete(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func (r roomMysqlRepository) SetRoomMap(ctx context.Context, key string, data room.RoomInfo) error {
	//TODO implement me
	panic("implement me")
}

func (r roomMysqlRepository) GetRoomMap(ctx context.Context, key, mapKey string) (room.RoomInfo, error) {
	//TODO implement me
	panic("implement me")
}
