package chat

import (
	"github.com/gorilla/websocket"
	"live-chat-server/domain/chat/types"
	"log/slog"
	"time"
)

type Client struct {
	Send   chan *message
	Room   *Room
	userId string
	Socket *websocket.Conn
}

func NewClient(socket *websocket.Conn, r *Room, clientId string) *Client {
	return &Client{
		Socket: socket,
		Send:   make(chan *message, types.MessageBufferSize),
		Room:   r,
		userId: clientId,
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	// client 가 메시지를 전송하는 함수

	for msg := range c.Send {
		if err := c.Socket.WriteJSON(msg); err != nil {
			slog.Error("failed write message, err : %v", err)
		}
	}
}

func (c *Client) Read() {
	defer c.Socket.Close()
	// client 가 메시지를 읽는 함수

	for {
		var msg *message
		if err := c.Socket.ReadJSON(&msg); err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				break
			}

			slog.Error("failed read message, err : %s", err.Error())
			continue
		}
		msg.Time = time.Now().Unix()
		msg.SendUserId = c.userId

		c.Room.Forward <- msg
	}
}
