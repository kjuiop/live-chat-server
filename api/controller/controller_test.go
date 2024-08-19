package controller

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

var (
	systemController *SystemController
	roomController   *RoomController
)

func TestMain(m *testing.M) {

	systemController = NewSystemController()

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
