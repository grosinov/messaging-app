package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name         string
		input        map[string]interface{}
		setupMock    func(mock *service.MockService)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "success",
			input: map[string]interface{}{
				"username": "testuser",
				"password": "testpass",
			},
			setupMock: func(mock *service.MockService) {
				mock.On("CreateUser", "testuser", "testpass").Return(&models.User{
					ID:       1,
					Username: "testuser",
					Password: "hashedpass",
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id": float64(1),
			},
		},
		{
			name: "failure - invalid request body",
			input: map[string]interface{}{
				"username": 1,
				"password": "",
			},
			setupMock:    func(mock *service.MockService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request body\n",
		},
		{
			name: "failure - invalid username",
			input: map[string]interface{}{
				"username": "",
				"password": "",
			},
			setupMock:    func(mock *service.MockService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid username\n",
		},
		{
			name: "failure - invalid password",
			input: map[string]interface{}{
				"username": "testuser",
				"password": "",
			},
			setupMock:    func(mock *service.MockService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid password\n",
		},
		{
			name: "failure - service error",
			input: map[string]interface{}{
				"username": "testuser",
				"password": "testpass",
			},
			setupMock: func(mock *service.MockService) {
				mock.
					On("CreateUser", "testuser", "testpass").
					Return(nil, errors.New("service error"))
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

			req := httptest.NewRequest(http.MethodGet, "/users", bytes.NewReader(jsonBytes))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

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
