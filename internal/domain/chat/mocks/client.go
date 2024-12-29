package mocks

import (
	"github.com/gorilla/websocket"
	"net"
	"time"
)

type MockWebSocketConn struct {
	MessagesSent    [][]byte // 보낸 메시지 저장
	MessagesToRead  [][]byte // 읽을 메시지 저장
	CloseWasCalled  bool     // Close 메서드가 호출되었는지 여부
	WriteMessageErr error    // WriteMessage 호출 시 발생할 에러
	ReadMessageErr  error    // ReadMessage 호출 시 발생할 에러
}

func (m *MockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	if m.WriteMessageErr != nil {
		return m.WriteMessageErr
	}
	m.MessagesSent = append(m.MessagesSent, data)
	return nil
}

func (m *MockWebSocketConn) ReadMessage() (int, []byte, error) {
	if m.ReadMessageErr != nil {
		return 0, nil, m.ReadMessageErr
	}

	if len(m.MessagesToRead) == 0 {
		return 0, nil, nil
	}

	message := m.MessagesToRead[0]
	m.MessagesToRead = m.MessagesToRead[1:]
	return websocket.TextMessage, message, nil
}

func (m *MockWebSocketConn) Close() error {
	m.CloseWasCalled = true
	return nil
}

func (m *MockWebSocketConn) LocalAddr() net.Addr                       { return nil }
func (m *MockWebSocketConn) RemoteAddr() net.Addr                      { return nil }
func (m *MockWebSocketConn) SetWriteDeadline(t time.Time) error        { return nil }
func (m *MockWebSocketConn) SetReadDeadline(t time.Time) error         { return nil }
func (m *MockWebSocketConn) SetPongHandler(h func(string) error) error { return nil }
