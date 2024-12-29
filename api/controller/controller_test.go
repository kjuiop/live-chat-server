package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/config"
	cu "live-chat-server/internal/domain/chat/usecase"
	"live-chat-server/internal/domain/room"
	rm "live-chat-server/internal/domain/room/usecase/mocks"
	sm "live-chat-server/internal/domain/system/usecase/mocks"
	"live-chat-server/internal/utils"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var testClient *TestClient

type TestClient struct {
	srv *httptest.Server

	systemController *SystemController
	roomController   *RoomController
	chatController   *ChatController
}

func TestMain(m *testing.M) {

	su := sm.NewSystemUseCaseStub()

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("fail to read config err : %v", err)
	}

	timeout := time.Duration(cfg.Policy.ContextTimeout) * time.Second

	data := utils.GetTestInitRoomData()
	var roomInfo []room.RoomInfo
	if err := json.Unmarshal(data, &roomInfo); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	us := rm.NewRoomUseCaseStub(roomInfo)
	chatUseCase := cu.NewChatUseCase(us, timeout)

	testClient = &TestClient{
		systemController: NewSystemController(su),
		roomController:   NewRoomController(cfg.Policy, us),
		chatController:   NewChatController(chatUseCase),
	}

	server := setupTestServer()
	testClient.srv = server

	exitCode := m.Run()
	server.Close()
	os.Exit(exitCode)
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

func setupTestServer() *httptest.Server {
	r := gin.Default()
	r.GET("/ws/chat/join/rooms/:room_id/user/:user_id", testClient.chatController.ServeWs)
	gin.SetMode(gin.TestMode)
	return httptest.NewServer(r)
}

func createWsURL(serverURL, testPath string) string {
	return fmt.Sprintf("ws%s%s", serverURL[len("http"):], testPath)
}
