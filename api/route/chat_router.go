package route

import (
	"github.com/gin-gonic/gin"
	"live-chat-server/api/controller"
	"live-chat-server/domain/chat"
)

func SetupChatGroup(ws *gin.RouterGroup, cu chat.ChatUseCase) {
	chatController := controller.NewChatController(cu)
	ws.GET("/chat/join/rooms/:room_id/user/:user_id", chatController.ServeWs)
}
