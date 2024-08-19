package controller

import (
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
