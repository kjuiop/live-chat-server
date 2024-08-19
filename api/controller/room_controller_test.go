package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"live-chat-server/domain"
	"live-chat-server/domain/room"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoomController_CreateChatRoom(t *testing.T) {

	tests := []struct {
		expectedCode int
		title        string
		request      room.RoomRequest
		response     domain.SuccessRes
	}{
		{
			expectedCode: http.StatusCreated, title: "Create Chat Room test success case",
			request: room.RoomRequest{
				CustomerId:   "jungin-kim",
				ChannelKey:   "calksdjfkdsa",
				BroadCastKey: "20240721-askdflj",
			},
			response: domain.SuccessRes{
				ErrorCode: domain.NoError,
				Message:   "ok",
				Result: room.RoomResponse{
					RoomId:       "N2-01J5MRM3M9AA83QT3RNHYRJ9RP",
					CustomerId:   "jungin-kim",
					ChannelKey:   "calksdjfkdsa",
					BroadcastKey: "20240721-askdflj",
					CreatedAt:    1724052541,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			testAssert := assert.New(t)
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)

			jsonRequest, err := json.Marshal(tc.request)
			if err != nil {
				t.Fatal(err)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/api/rooms", bytes.NewBuffer(jsonRequest))

			roomController.CreateChatRoom(c)

			testAssert.Equal(tc.expectedCode, resp.Code)

			var responseBody domain.SuccessRes
			if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			testAssert.NoError(err)

			testAssert.Equal(tc.response.ErrorCode, responseBody.ErrorCode)
			testAssert.Equal(tc.response.Message, responseBody.Message)
		})
	}
}
