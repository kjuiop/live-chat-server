package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/domain"
	"net/http"
)

type SystemController struct {
}

func NewSystemController() *SystemController {
	return &SystemController{}
}

func (s *SystemController) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, domain.HealthRes{Message: domain.GetCustomMessage(domain.NoError)})
}

func (s *SystemController) OccurPanic(c *gin.Context) {

	requestId, exists := c.Get("request_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, domain.FailRes{Message: "request not exist"})
		return
	}

	panic(fmt.Errorf("panic encounter, request_id : %s", requestId))
}
