package controller

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/config"
	"live-chat-server/internal/domain/room"
	"live-chat-server/internal/domain/room/usecase/mocks"
	"log"
	"os"
	"testing"
)

var (
	systemController *SystemController
	roomController   *RoomController
)

//go:embed testdata/rooms/init_data.json
var initRooms embed.FS

func TestMain(m *testing.M) {

	systemController = NewSystemController()

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v", err)
	}

	data, err := initRooms.ReadFile("testdata/rooms/init_data.json")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	var roomInfo []room.RoomInfo
	if err := json.Unmarshal(data, &roomInfo); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	us := mocks.NewRoomUseCaseStub(roomInfo)
	roomController = NewRoomController(cfg.Policy, us)

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func convertToBytes[T any](data T) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	return jsonData, nil
}

func convertResultTo[T any](result interface{}, v *T) error {
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(resultBytes, v); err != nil {
		return err
	}
	return nil
}
