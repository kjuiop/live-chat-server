package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/internal/domain"
	"live-chat-server/internal/domain/system"
	"net/http"
)

type SystemController struct {
	SystemUseCase system.UseCase
}

func NewSystemController(useCase system.UseCase) *SystemController {
	return &SystemController{
		SystemUseCase: useCase,
	}
}

func (s *SystemController) successResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, domain.ApiResponse{
		ErrorCode: domain.NoError,
		Message:   domain.GetCustomMessage(domain.NoError),
		Result:    data,
	})
}

func (s *SystemController) GetHealth(c *gin.Context) {
	s.successResponse(c, http.StatusOK, nil)
	return
}

func (s *SystemController) GetServerList(c *gin.Context) {
	s.successResponse(c, http.StatusOK, nil)
	return
}

func (s *SystemController) OccurPanic(c *gin.Context) {

	requestId, exists := c.Get("request_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, domain.ApiResponse{Message: "request not exist"})
		return
	}

	panic(fmt.Errorf("panic encounter, request_id : %s", requestId))
}
