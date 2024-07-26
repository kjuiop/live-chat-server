package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/config"
	"live-chat-server/models"
	"net/http"
)

type RoomController struct {
	cfg         config.RoomPolicy
	RoomUseCase models.RoomUseCase
}

func NewRoomController(cfg config.RoomPolicy, useCase models.RoomUseCase) *RoomController {
	return &RoomController{
		cfg:         cfg,
		RoomUseCase: useCase,
	}
}

func (r *RoomController) successResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, models.SuccessRes{
		ErrorCode: models.NoError,
		Message:   models.GetCustomMessage(models.NoError),
		Result:    data,
	})
}

func (r *RoomController) failResponse(c *gin.Context, statusCode, errorCode int, err error) {

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

func (r *RoomController) CreateChatRoom(c *gin.Context) {
	req := models.CreateRoomReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		r.failResponse(c, http.StatusBadRequest, models.ErrParsing, fmt.Errorf("CreateRoom json parsing err : %w", err))
		return
	}

	roomInfo := models.NewRoomInfo(&req, r.cfg.Prefix)
	if err := r.RoomUseCase.CreateChatRoom(c, roomInfo); err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMSETError, fmt.Errorf("CreateRoom HMSET err : %w", err))
		return
	}

	r.successResponse(c, http.StatusCreated, roomInfo)
}

func (r *RoomController) GetChatRoom(c *gin.Context) {

	roomId := c.Param("roomId")
	roomInfo, err := r.RoomUseCase.GetChatRoomById(c, roomId)
	if err != nil {
		r.failResponse(c, http.StatusNotFound, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, err : %w", err))
		return
	}

	r.successResponse(c, http.StatusOK, roomInfo)
}
