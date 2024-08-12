package route

import (
	"github.com/gin-gonic/gin"
	"live-chat-server/config"
	redis "live-chat-server/internal/redis"
	"live-chat-server/repository"
	"live-chat-server/usecase"
	"time"
)

func Setup(api, ws *gin.RouterGroup, cfg config.Policy, timeout time.Duration, redis redis.Client) {

	SetupSystemGroup(api)

	rr := repository.NewRoomRepository(redis)
	ur := usecase.NewRoomUseCase(rr, timeout)

	setupRoomGroup(api, cfg, ur)

	cu := usecase.NewChatUseCase(ur, timeout)
	setupChatGroup(ws, cu)
}
