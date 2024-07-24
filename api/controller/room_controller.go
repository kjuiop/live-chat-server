package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/config"
	"live-chat-server/models"
	"log/slog"
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

	message := models.GetCustomErrMessage(errorCode, err.Error())
	slog.Error(message, "status_code", statusCode, "error_code", errorCode, "err", err)

	c.JSON(statusCode, models.FailRes{
		ErrorCode: errorCode,
		Message:   message,
	})
}

func (r *RoomController) CreateRoom(c *gin.Context) {
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
