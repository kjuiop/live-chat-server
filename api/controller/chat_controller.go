package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"live-chat-server/chat"
	"live-chat-server/chat/types"
	"live-chat-server/models"
	"net/http"
	"sync"
)

var crMutex = &sync.RWMutex{}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  types.SocketBufferSize,
	WriteBufferSize: types.MessageBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatController struct {
	RoomUseCase models.RoomUseCase
	hub         map[string]*chat.Room
}

func NewChatController(roomUseCase models.RoomUseCase) *ChatController {
	return &ChatController{
		RoomUseCase: roomUseCase,
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

func (cc *ChatController) JoinChatRoom(c *gin.Context) {

	roomId := c.Param("room_id")
	userId := c.Param("user_id")

	chatRoom, err := cc.getChatRoom(c, roomId)
	if err != nil {
		cc.failResponse(c, http.StatusNotFound, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, roomId : %s, err : %w", roomId, err))
		return
	}

	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		cc.failResponse(c, http.StatusConflict, models.ErrNotConnectSocket, fmt.Errorf("failed connect socket, roomId : %s, err : %w", roomId, err))
		return
	}

	client := chat.NewClient(socket, chatRoom, userId)

	chatRoom.Join <- client

	defer func() {
		chatRoom.Leave <- client
	}()

	go client.Write()

	client.Read()

	cc.successResponse(c, http.StatusOK, models.SuccessRes{
		ErrorCode: 0,
		Message:   "close chat connection",
		Result:    nil,
	})
}

func (cc *ChatController) getChatRoom(c *gin.Context, roomId string) (*chat.Room, error) {

	crMutex.Lock()
	defer func() {
		crMutex.Unlock()
	}()

	if _, ok := cc.hub[roomId]; !ok {
		roomInfo, err := cc.RoomUseCase.GetChatRoomById(c, roomId)
		if err != nil {
			return nil, fmt.Errorf("not found chat room, key : %s, err : %w", roomId, err)
		}
		cc.hub[roomId] = chat.NewChatRoom(roomInfo)
	}

	return cc.hub[roomId], nil
}
