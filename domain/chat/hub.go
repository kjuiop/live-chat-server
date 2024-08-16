package chat

import (
	"live-chat-server/domain/room"
)

type Room struct {
	RoomId string `json:"RoomId"`
	Alive  bool   `json:"Alive"`

	Forward chan *message // 수신되는 메시지를 보관하는 값
	// 들어오는 메시지를 다른 클라이언트에게 전송을 합니다.

	Join  chan *Client // Socket 이 연결되는 경우에 적용
	Leave chan *Client // Socket 이 끊어지는 경우에 대해서 적용

	Clients map[*Client]bool // 현재 방에 있는 Client 정보를 저장
}

func NewChatRoom(roomInfo room.RoomInfo) *Room {
	chatRoom := &Room{
		RoomId:  roomInfo.RoomId,
		Alive:   true,
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
	go chatRoom.chatInit()
	return chatRoom
}

func (r *Room) chatInit() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:
			// msg type 에 따른 분기 처리 가능
			switch msg.Method {
			case "chat":
				r.broadcastChat(msg)
			}
		}
	}
}

func (r *Room) broadcastChat(msg *message) {
	for client := range r.Clients {
		client.Send <- msg
	}
}
