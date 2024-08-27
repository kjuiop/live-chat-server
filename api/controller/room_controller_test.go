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
		expectedResp domain.ApiResponse
	}{
		{
			expectedCode: http.StatusCreated, title: "Create Chat Room test success case",
			request: room.RoomRequest{
				CustomerId:   "jungin-kim",
				ChannelKey:   "calksdjfkdsa",
				BroadCastKey: "20240721-askdflj",
			},
			expectedResp: domain.ApiResponse{
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
		{
			expectedCode: http.StatusBadRequest, title: "Create Chat Room Bad Request Fail case",
			request: room.RoomRequest{
				CustomerId:   "jungin-kim",
				ChannelKey:   "",
				BroadCastKey: "20240721-askdflj",
			},
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.ErrParsing,
				Message:   "invalid request body",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			testAssert := assert.New(t)
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)

			jsonRequest, err := convertToBytes(tc.request)
			if err != nil {
				t.Fatal(err)
			}
			c.Request = httptest.NewRequest(http.MethodPost, "/api/rooms", bytes.NewBuffer(jsonRequest))

			roomController.CreateChatRoom(c)

			testAssert.Equal(tc.expectedCode, resp.Code)

			var responseBody *domain.ApiResponse
			if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			testAssert.Equal(tc.expectedResp.ErrorCode, responseBody.ErrorCode)
			testAssert.Equal(tc.expectedResp.Message, responseBody.Message)

			if tc.expectedCode == http.StatusCreated {

				testAssert.NotNil(responseBody.Result)

				var roomResp room.RoomResponse
				if err := convertResultTo(responseBody.Result, &roomResp); err != nil {
					t.Fatal(err)
				}

				testAssert.NotEmpty(roomResp.RoomId, "room_id is not empty")
				testAssert.True(roomResp.CreatedAt > 0, "created_at should be large to 0")

				testAssert.Equal(tc.expectedResp.Result.(room.RoomResponse).CustomerId, roomResp.CustomerId)
				testAssert.Equal(tc.expectedResp.Result.(room.RoomResponse).ChannelKey, roomResp.ChannelKey)
				testAssert.Equal(tc.expectedResp.Result.(room.RoomResponse).BroadcastKey, roomResp.BroadcastKey)
			}
		})
	}
}
