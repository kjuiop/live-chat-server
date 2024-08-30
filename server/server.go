package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"sync"
)

type Client interface {
	Run(wg *sync.WaitGroup)
	Shutdown(ctx context.Context)
	GetRouterGroup(name string) *gin.RouterGroup
}
