package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/config"
	"live-chat-server/models"
	"net/http"
)

type RoomController struct {
	cfg         config.Policy
	RoomUseCase models.RoomUseCase
}

func NewRoomController(cfg config.Policy, useCase models.RoomUseCase) *RoomController {
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
	req := models.RoomRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		r.failResponse(c, http.StatusBadRequest, models.ErrParsing, fmt.Errorf("CreateRoom json parsing err : %w", err))
		return
	}

	roomInfo := models.NewRoomInfo(&req, r.cfg.Prefix)
	if err := r.RoomUseCase.CreateChatRoom(c, roomInfo); err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMSETError, fmt.Errorf("CreateRoom HMSET err : %w", err))
		return
	}

	if err := r.RoomUseCase.RegisterRoomId(c, roomInfo); err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMSETError, fmt.Errorf("register room id HMSET err : %w", err))
		return
	}

	roomRes := models.RoomResponse{
		RoomId:       roomInfo.RoomId,
		CustomerId:   roomInfo.CustomerId,
		ChannelKey:   roomInfo.ChannelKey,
		BroadcastKey: roomInfo.BroadcastKey,
		CreatedAt:    roomInfo.CreatedAt,
	}

	r.successResponse(c, http.StatusCreated, roomRes)
}

func (r *RoomController) GetChatRoomId(c *gin.Context) {
	req := models.RoomIdRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		r.failResponse(c, http.StatusBadRequest, models.ErrParsing, fmt.Errorf("GetChatRoomId bind fail, err : %w", err))
		return
	}

	roomInfo, err := r.RoomUseCase.GetChatRoomId(c, req)
	if err != nil {
		r.failResponse(c, http.StatusNotFound, models.ErrNotFoundChatRoom, fmt.Errorf("GetChatRoomId fail get roomInfo, err : %w", err))
		return
	}

	roomRes := models.RoomResponse{
		RoomId: roomInfo.RoomId,
	}

	r.successResponse(c, http.StatusOK, roomRes)
}

func (r *RoomController) GetChatRoom(c *gin.Context) {

	roomId := c.Param("roomId")
	roomInfo, err := r.RoomUseCase.GetChatRoomById(c, roomId)
	if err != nil {
		r.failResponse(c, http.StatusNotFound, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, err : %w", err))
		return
	}

	roomRes := models.RoomResponse{
		RoomId:       roomInfo.RoomId,
		CustomerId:   roomInfo.CustomerId,
		ChannelKey:   roomInfo.ChannelKey,
		BroadcastKey: roomInfo.BroadcastKey,
		CreatedAt:    roomInfo.CreatedAt,
	}

	r.successResponse(c, http.StatusOK, roomRes)
}

func (r *RoomController) UpdateChatRoom(c *gin.Context) {

	roomId := c.Param("roomId")
	req := models.RoomRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		r.failResponse(c, http.StatusBadRequest, models.ErrParsing, fmt.Errorf("UpdateChatRoom json parsing err : %w", err))
		return
	}

	isExist, err := r.RoomUseCase.CheckExistRoomId(c, roomId)
	if err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisExistError, fmt.Errorf("fail exec redis exist cmd, err : %w", err))
		return
	}

	if !isExist {
		r.failResponse(c, http.StatusNotFound, models.ErrNotFoundChatRoom, fmt.Errorf("not found chat room, err : %w", err))
		return
	}

	roomInfo := models.UpdateRoomInfo(&req, roomId)
	savedInfo, err := r.RoomUseCase.UpdateChatRoom(c, roomId, roomInfo)
	if err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMSETError, fmt.Errorf("fail exec redis save cmd, err : %w", err))
		return
	}

	roomRes := models.RoomResponse{
		RoomId:       savedInfo.RoomId,
		CustomerId:   savedInfo.CustomerId,
		ChannelKey:   savedInfo.ChannelKey,
		BroadcastKey: savedInfo.BroadcastKey,
		CreatedAt:    savedInfo.CreatedAt,
	}

	r.successResponse(c, http.StatusOK, roomRes)
}

func (r *RoomController) DeleteChatRoom(c *gin.Context) {

	roomId := c.Param("roomId")
	isExist, err := r.RoomUseCase.CheckExistRoomId(c, roomId)
	if err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisExistError, fmt.Errorf("fail exec redis exist cmd, err : %w", err))
		return
	}

	if !isExist {
		r.successResponse(c, http.StatusNoContent, nil)
		return
	}

	if err := r.RoomUseCase.DeleteChatRoom(c, roomId); err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMDELError, fmt.Errorf("fail exec redis hmdel cmd, err : %w", err))
		return
	}

	r.successResponse(c, http.StatusOK, nil)
}
