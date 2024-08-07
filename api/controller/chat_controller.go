package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/models"
	"net/http"
)

type ChatController struct {
	RoomUseCase models.RoomUseCase
}

func NewChatController(roomUseCase models.RoomUseCase) *ChatController {
	return &ChatController{
		RoomUseCase: roomUseCase,
	}
}

func (cc *ChatController) successResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, models.SuccessRes{
		ErrorCode: models.NoError,
		Message:   models.GetCustomMessage(models.NoError),
		Result:    data,
	})
}

func (cc *ChatController) failResponse(c *gin.Context, statusCode, errorCode int, err error) {

	logMessage := models.GetCustomErrMessage(errorCode, err.Error())
	c.Errors = append(c.Errors, &gin.Error{
		Err:  fmt.Errorf(logMessage),
		Type: gin.ErrorTypePrivate,
	})

	c.JSON(statusCode, models.FailRes{
		ErrorCode: errorCode,
		Message:   models.GetCustomMessage(errorCode),
	})
}

func (cc *ChatController) JoinChatRoom(c *gin.Context) {

	roomId := c.Param("roomId")

	_, err := cc.getChatRoom(c, roomId)
	if err != nil {
		cc.failResponse(c, http.StatusNotFound, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, roomId : %s, err : %w", roomId, err))
		return
	}

	cc.successResponse(c, http.StatusOK, models.SuccessRes{
		ErrorCode: 0,
		Message:   "close chat connection",
		Result:    nil,
	})
}

func (cc *ChatController) getChatRoom(c *gin.Context, roomId string) (*models.ChatRoom, error) {

	roomInfo, err := cc.RoomUseCase.GetChatRoomById(c, roomId)
	if err != nil {
		return nil, fmt.Errorf("not found chat room, key : %s, err : %w", roomId, err)
	}

	chatRoom := models.NewChatRoom(roomInfo)
	return chatRoom, nil
}
