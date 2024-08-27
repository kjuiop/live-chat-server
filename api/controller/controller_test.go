package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/config"
	"live-chat-server/usecase/mocks"
	"log"
	"os"
	"testing"
)

var (
	systemController *SystemController
	roomController   *RoomController
)

func TestMain(m *testing.M) {

	systemController = NewSystemController()

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v\n", err)
	}

	us := mocks.NewRoomUseCaseStub()
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
