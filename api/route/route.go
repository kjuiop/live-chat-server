package route

import (
	"github.com/gin-gonic/gin"
	"live-chat-server/config"
	"live-chat-server/database"
	"live-chat-server/repository"
	"live-chat-server/usecase"
	"time"
)

func Setup(api, ws *gin.RouterGroup, cfg config.Policy, timeout time.Duration, db database.Client) {

	SetupSystemGroup(api)

	rr := repository.NewRoomRepository(db)
	ur := usecase.NewRoomUseCase(rr, timeout)

	setupRoomGroup(api, cfg, ur)

	cu := usecase.NewChatUseCase(ur, timeout)
	setupChatGroup(ws, cu)
}
