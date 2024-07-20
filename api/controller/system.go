package controller

import (
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
	c.JSON(http.StatusOK, models.HealthRes{Message: "pong"})
}
