package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	timestamp := time.DateTime
	tests := []struct {
		name         string
		requestUser  uint64
		input        map[string]interface{}
		setupMock    func(mock *service.MockService)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:        "success",
			requestUser: 1,
			input: map[string]interface{}{
				"sender":    1,
				"recipient": 2,
				"content": map[string]interface{}{
					"type": "text",
					"text": "Hello",
				},
			},
			setupMock: func(mock *service.MockService) {
				mock.On("GetUser", uint64(1)).Return(nil, nil)
				mock.On("GetUser", uint64(2)).Return(nil, nil)
				mock.On("SendMessage", uint64(1), uint64(2), &models.Content{
					Type: "text",
					Text: "Hello",
				}).Return(&models.Message{
					SenderID:    1,
					RecipientID: 2,
					Content: models.Content{
						Type: "text",
						Text: "Hello",
					},
					Timestamp: timestamp,
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":        0.0,
				"timestamp": "2006-01-02 15:04:05",
			},
		},
		{
			name:        "failure - invalid sender",
			requestUser: 1,
			input: map[string]interface{}{
				"sender":    1,
				"recipient": 2,
				"content": map[string]interface{}{
					"type": "text",
					"text": "Hello",
				},
			},
			setupMock: func(mock *service.MockService) {
				mock.On("GetUser", uint64(1)).Return(nil, httperrors.BadRequestError("user not found"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "user not found\n",
		},
		{
			name:        "failure - invalid recipient",
			requestUser: 1,
			input: map[string]interface{}{
				"sender":    1,
				"recipient": 2,
				"content": map[string]interface{}{
					"type": "text",
					"text": "Hello",
				},
			},
			setupMock: func(mock *service.MockService) {
				mock.On("GetUser", uint64(0)).Return(nil, nil)
				mock.On("GetUser", uint64(1)).Return(nil, httperrors.BadRequestError("user not found"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "user not found\n",
		},
		{
			name: "failure - invalid request body",
			input: map[string]interface{}{
				"sender": "invalid",
			},
			setupMock:    func(mock *service.MockService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request body\n",
		},
		{
			name: "failure - service error",
			input: map[string]interface{}{
				"sender":    uint64(1),
				"recipient": uint64(2),
				"content": map[string]interface{}{
					"type": "text",
					"text": "Hello",
				},
			},
			setupMock: func(mock *service.MockService) {
				mock.On("SendMessage", uint64(1), uint64(2), &models.Content{
					Type: "text",
					Text: "Hello",
				}).Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal Server Error: service error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(service.MockService)
			tt.setupMock(mockService)
			handler := NewHandler(mockService)

			jsonBytes, _ := json.Marshal(tt.input)

			req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(jsonBytes))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", uint64(tt.requestUser)))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.zCGBEiC4n4X5jij4lK4nSEtrbebYxELZ6OfBwdm6CJg")
			w := httptest.NewRecorder()

			handler.SendMessage(w, req)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				response := w.Body.String()
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
