package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"live-chat-server/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealthSuccess(t *testing.T) {

	tests := []struct {
		expectedCode int
		title        string
		response     domain.ApiResponse
	}{
		{
			http.StatusOK, "health-check api test", domain.ApiResponse{
				ErrorCode: domain.NoError,
				Message:   domain.GetCustomMessage(domain.NoError),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			testAssert := assert.New(t)
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)

			c.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)

			systemController.GetHealth(c)
			testAssert.Equal(tc.expectedCode, resp.Code)

			var responseBody domain.ApiResponse
			err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
			testAssert.NoError(err)
			testAssert.Equal(tc.response.ErrorCode, responseBody.ErrorCode)
			testAssert.Equal(tc.response.Message, responseBody.Message)
		})
	}
}
