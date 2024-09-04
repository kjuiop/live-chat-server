package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"live-chat-server/internal/reporter"
	"log/slog"
	"net/http"
)

func RecoveryErrorReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("recovered from panic : %v", err)
				slog.Error(errMsg)
				reporter.Client.SendSlackPanicReport(errMsg)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
