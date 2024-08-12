package route

import (
	"github.com/gin-gonic/gin"
	"live-chat-server/api/controller"
	"live-chat-server/config"
	"live-chat-server/models"
)

func setupRoomGroup(api *gin.RouterGroup, cfg config.Policy, ur models.RoomUseCase) {

	roomController := controller.NewRoomController(cfg, ur)
	api.POST("/rooms", roomController.CreateChatRoom)
	api.GET("/rooms/:room_id", roomController.GetChatRoom)
	api.PUT("/rooms/:room_id", roomController.UpdateChatRoom)
	api.DELETE("/rooms/:room_id", roomController.DeleteChatRoom)
	api.GET("/rooms/id", roomController.GetChatRoomId)
}
