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

func TestRoomController_CreateChatRoom_Success(t *testing.T) {

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

			// response body 검증
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

func TestRoomController_CreateChatRoom_Fail(t *testing.T) {

	tests := []struct {
		expectedCode int
		title        string
		request      map[string]interface{}
		expectedResp domain.ApiResponse
	}{
		{
			expectedCode: http.StatusBadRequest, title: "empty request param",
			request: map[string]interface{}{},
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.ErrParsing,
				Message:   "invalid request body",
			},
		},
		{
			expectedCode: http.StatusBadRequest, title: "invalid data type field",
			request: map[string]interface{}{"customer_id": 1, "channel_key": 1, "broadcast_key": 1},
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
		})
	}
}

func TestRoomController_GetChatRoom_Success(t *testing.T) {

	tests := []struct {
		expectedCode int
		title        string
		roomId       string
		expectedResp domain.ApiResponse
	}{
		{
			expectedCode: http.StatusOK, title: "Get Chat Room test success case",
			roomId: "N1-TESTMRM3M9AA83QT3RNHYRJ9RP",
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.NoError,
				Message:   "ok",
				Result: room.RoomResponse{
					RoomId:       "N1-TESTMRM3M9AA83QT3RNHYRJ9RP",
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

			c.Request = httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
			c.Params = gin.Params{
				{Key: "room_id", Value: tc.roomId},
			}
			roomController.GetChatRoom(c)

			testAssert.Equal(tc.expectedCode, resp.Code)

			var responseBody *domain.ApiResponse
			if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			testAssert.Equal(tc.expectedResp.ErrorCode, responseBody.ErrorCode)
			testAssert.Equal(tc.expectedResp.Message, responseBody.Message)

			// response body 검증
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
		})
	}
}

func TestRoomController_GetChatRoom_Fail(t *testing.T) {

	tests := []struct {
		expectedCode int
		title        string
		roomId       string
		expectedResp domain.ApiResponse
	}{
		{
			expectedCode: http.StatusBadRequest, title: "Get Chat Room test fail case, empty room_id",
			roomId: "",
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.ErrEmptyParam,
				Message:   "invalid params",
			},
		},
		{
			expectedCode: http.StatusNotFound, title: "Get Chat Room test fail case, not exist",
			roomId: "N1-NOT_EXIST_ROOM_ID",
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.ErrNotFoundChatRoom,
				Message:   "not found chat room",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			testAssert := assert.New(t)
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
			c.Params = gin.Params{
				{Key: "room_id", Value: tc.roomId},
			}
			roomController.GetChatRoom(c)

			testAssert.Equal(tc.expectedCode, resp.Code)

			var responseBody *domain.ApiResponse
			if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			testAssert.Equal(tc.expectedResp.ErrorCode, responseBody.ErrorCode)
			testAssert.Equal(tc.expectedResp.Message, responseBody.Message)
		})
	}
}

func TestRoomController_UpdateChatRoom_Success(t *testing.T) {

	tests := []struct {
		expectedCode int
		title        string
		roomId       string
		request      room.RoomRequest
		expectedResp domain.ApiResponse
	}{
		{
			expectedCode: http.StatusOK, title: "Update Chat Room test success case",
			roomId: "N1-TESTMRM3M9AA83QT3RNHYRJ9RP",
			request: room.RoomRequest{
				CustomerId:   "test",
				ChannelKey:   "test",
				BroadCastKey: "test",
			},
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.NoError,
				Message:   "ok",
				Result: room.RoomResponse{
					RoomId:       "N1-TESTMRM3M9AA83QT3RNHYRJ9RP",
					CustomerId:   "test",
					ChannelKey:   "test",
					BroadcastKey: "test",
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

			jsonRequest, err := convertToBytes(tc.request)
			if err != nil {
				t.Fatal(err)
			}
			c.Request = httptest.NewRequest(http.MethodPut, "/api/rooms", bytes.NewBuffer(jsonRequest))
			c.Params = gin.Params{
				{Key: "room_id", Value: tc.roomId},
			}

			roomController.UpdateChatRoom(c)

			testAssert.Equal(tc.expectedCode, resp.Code)

			var responseBody *domain.ApiResponse
			if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			testAssert.Equal(tc.expectedResp.ErrorCode, responseBody.ErrorCode)
			testAssert.Equal(tc.expectedResp.Message, responseBody.Message)

			// response body 검증
			if tc.expectedCode == http.StatusOK {

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

func TestRoomController_UpdateChatRoom_Fail(t *testing.T) {

	tests := []struct {
		expectedCode int
		title        string
		roomId       string
		request      map[string]interface{}
		expectedResp domain.ApiResponse
	}{
		{
			expectedCode: http.StatusBadRequest, title: "Update Chat Room empty param",
			roomId:  "",
			request: map[string]interface{}{},
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.ErrEmptyParam,
				Message:   "invalid params",
			},
		},
		{
			expectedCode: http.StatusBadRequest, title: "Update Chat Room empty request body",
			roomId:  "N1-TESTMRM3M9AA83QT3RNHYRJ9RP",
			request: map[string]interface{}{},
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.ErrParsing,
				Message:   "invalid request body",
			},
		},
		{
			expectedCode: http.StatusNotFound, title: "Update Chat Room not exist room_id",
			roomId:  "N1-NOT_EXIST_MRM3M9AA83QT3RNHYRJ9RP",
			request: map[string]interface{}{"customer_id": "jungin-kim", "channel_key": "calksdjfkdsa", "broadcast_key": "20240721-askdflj"},
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.ErrNotFoundChatRoom,
				Message:   "not found chat room",
			},
		},
		{
			expectedCode: http.StatusBadRequest, title: "Update Chat Room invalid request body",
			roomId:  "N1-TESTMRM3M9AA83QT3RNHYRJ9RP",
			request: map[string]interface{}{"customer_id": 1, "channel_key": 1, "broadcast_key": 1},
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
			c.Request = httptest.NewRequest(http.MethodPut, "/api/rooms", bytes.NewBuffer(jsonRequest))
			c.Params = gin.Params{
				{Key: "room_id", Value: tc.roomId},
			}

			roomController.UpdateChatRoom(c)

			testAssert.Equal(tc.expectedCode, resp.Code)

			var responseBody *domain.ApiResponse
			if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			testAssert.Equal(tc.expectedResp.ErrorCode, responseBody.ErrorCode)
			testAssert.Equal(tc.expectedResp.Message, responseBody.Message)
		})
	}
}

func TestRoomController_DeleteChatRoom_Success(t *testing.T) {

	tests := []struct {
		expectedCode int
		title        string
		roomId       string
		expectedResp domain.ApiResponse
	}{
		{
			expectedCode: http.StatusOK, title: "Delete Chat Room test success case",
			roomId: "N1-TESTMRM3M9AA83QT3RNHYRJ9RP",
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.NoError,
				Message:   "ok",
			},
		},
		{
			expectedCode: http.StatusNoContent, title: "Delete Chat Room test success case, No Content",
			roomId: "N1-NOT_EXTIST_3M9AA83QT3RNHYRJ9RP",
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.NoError,
				Message:   "ok",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			testAssert := assert.New(t)
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)

			c.Request = httptest.NewRequest(http.MethodDelete, "/api/rooms", nil)
			c.Params = gin.Params{
				{Key: "room_id", Value: tc.roomId},
			}

			roomController.DeleteChatRoom(c)

			testAssert.Equal(tc.expectedCode, resp.Code)

			if tc.expectedCode == http.StatusOK {
				var responseBody *domain.ApiResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
					t.Fatal(err)
				}

				testAssert.Equal(tc.expectedResp.ErrorCode, responseBody.ErrorCode)
				testAssert.Equal(tc.expectedResp.Message, responseBody.Message)
			}
		})
	}
}

func TestRoomController_DeleteChatRoom_Fail(t *testing.T) {
	tests := []struct {
		expectedCode int
		title        string
		roomId       string
		request      map[string]interface{}
		expectedResp domain.ApiResponse
	}{
		{
			expectedCode: http.StatusBadRequest, title: "Update Chat Room empty param",
			roomId:  "",
			request: map[string]interface{}{},
			expectedResp: domain.ApiResponse{
				ErrorCode: domain.ErrEmptyParam,
				Message:   "invalid params",
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
			c.Request = httptest.NewRequest(http.MethodDelete, "/api/rooms", bytes.NewBuffer(jsonRequest))
			c.Params = gin.Params{
				{Key: "room_id", Value: tc.roomId},
			}

			roomController.DeleteChatRoom(c)

			testAssert.Equal(tc.expectedCode, resp.Code)

			var responseBody *domain.ApiResponse
			if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			testAssert.Equal(tc.expectedResp.ErrorCode, responseBody.ErrorCode)
			testAssert.Equal(tc.expectedResp.Message, responseBody.Message)
		})
	}
}
