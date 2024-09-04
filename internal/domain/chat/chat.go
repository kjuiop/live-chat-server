package chat

import (
	"context"
	"net/http"
)

type ChatUseCase interface {
	ServeWs(ctx context.Context, writer http.ResponseWriter, request *http.Request, roomId, userId string) error
}
