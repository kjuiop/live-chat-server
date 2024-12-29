package controller

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"live-chat-server/internal/domain/chat"
	"net/http"
	"testing"
	"time"
)

func TestChatController_Websocket_Connect(t *testing.T) {

	tests := []struct {
		title       string
		wsURL       string
		isConnected bool
	}{
		{
			isConnected: true,
			title:       "ws connect success",
			wsURL:       "/ws/chat/join/rooms/N1-TESTMRM3M9AA83QT3RNHYRJ9RP/user/user1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {

			wsURL := createWsURL(testClient.srv.URL, tc.wsURL)
			conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if tc.isConnected {
				assert.NoError(t, err, "expected connection to succeed, but got error: %v", err)
				assert.NotNil(t, conn, "expected webSocket connection, but got nil")
				assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "expected status code 101 for websocket upgrade")
			} else {
				assert.Error(t, err, "expected connection to fail, but it succeeded")
				assert.Nil(t, conn, "expected nil connection on failure, but got a valid connection")
				if resp != nil {
					assert.NotEqual(t, http.StatusSwitchingProtocols, resp.StatusCode, "not expect status code 101 for failed connection")
				}
			}

			if conn != nil {
				conn.Close()
			}
		})
	}
}

func TestChatController_Ping_Pong(t *testing.T) {

	tests := []struct {
		title    string
		wsURL    string
		request  *chat.Message
		expected *chat.Message
	}{
		{
			title: "ws connect success",
			wsURL: "/ws/chat/join/rooms/N1-TESTMRM3M9AA83QT3RNHYRJ9RP/user/user1",
			request: &chat.Message{
				SendUserId: "user1",
				Message:    "hello",
				Time:       time.Now().Unix(),
				Method:     "chat",
			},
			expected: &chat.Message{
				SendUserId: "user1",
				Message:    "hello",
				Method:     "chat",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {

			wsURL := createWsURL(testClient.srv.URL, tc.wsURL)
			conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
			assert.NoError(t, err, "expected connection to succeed, but got error: %v", err)
			assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

			err = conn.WriteJSON(tc.request)
			assert.NoError(t, err)

			_, readData, err := conn.ReadMessage()
			assert.NoError(t, err)

			var respMsg chat.Message
			if err := json.Unmarshal(readData, &respMsg); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.expected.Message, respMsg.Message)
			assert.Equal(t, tc.expected.SendUserId, respMsg.SendUserId)
			assert.Equal(t, tc.expected.Method, respMsg.Method)

			time.Sleep(100 * time.Millisecond)
		})
	}
}
