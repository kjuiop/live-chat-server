package server

import (
	"context"
	"sync"
)

type Client interface {
	Run(wg *sync.WaitGroup)
	Shutdown(ctx context.Context)
}
