package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/chat"
	"live-chat-server/models"
	"net/http"
)

type ChatController struct {
	ChatUseCase models.ChatUseCase
	hub         map[string]*chat.Room
}

func NewChatController(chatUseCase models.ChatUseCase) *ChatController {
	return &ChatController{
		ChatUseCase: chatUseCase,
		hub:         make(map[string]*chat.Room),
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

func (cc *ChatController) ServeWs(c *gin.Context) {

	roomId := c.Param("room_id")
	userId := c.Param("user_id")

	if err := cc.ChatUseCase.ServeWs(c, c.Writer, c.Request, roomId, userId); err != nil {
		cc.failResponse(c, http.StatusBadRequest, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, roomId : %s, err : %w", roomId, err))
		return
	}

	cc.successResponse(c, http.StatusOK, models.SuccessRes{
		ErrorCode: 0,
		Message:   "close chat connection",
		Result:    nil,
	})
}
