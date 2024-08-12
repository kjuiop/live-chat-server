package route

import (
	"github.com/gin-gonic/gin"
	"live-chat-server/api/controller"
)

func SetupSystemGroup(router *gin.RouterGroup) {
	systemController := controller.NewSystemController()
	router.GET("/health-check", systemController.GetHealth)
	router.GET("/panic-test", systemController.OccurPanic)
}
