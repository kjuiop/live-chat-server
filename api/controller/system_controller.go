package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/models"
	"net/http"
)

type SystemController struct {
}

func NewSystemController() *SystemController {
	return &SystemController{}
}

func (s *SystemController) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthRes{Message: models.GetCustomMessage(models.NoError)})
}

func (s *SystemController) OccurPanic(c *gin.Context) {

	requestId, exists := c.Get("request_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, models.FailRes{Message: "request not exist"})
		return
	}

	panic(fmt.Errorf("panic encounter, request_id : %s", requestId))
}
