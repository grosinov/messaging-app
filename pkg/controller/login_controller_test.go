package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
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
				mock.On("Login", "testuser", "testpass").Return(uint64(1), "token123", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":    1.0,
				"token": "token123",
			},
		},
		{
			name: "failure - invalid request body",
			input: map[string]interface{}{
				"username": 1,
				"password": "testpass",
			},
			setupMock:    func(mock *service.MockService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request body\n",
		},
		{
			name: "failure - empty username",
			input: map[string]interface{}{
				"username": "",
				"password": "testpass",
			},
			setupMock:    func(mock *service.MockService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid username or password\n",
		},
		{
			name: "failure - empty password",
			input: map[string]interface{}{
				"username": "testuser",
				"password": "",
			},
			setupMock:    func(mock *service.MockService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid username or password\n",
		},
		{
			name: "failure - invalid credentials",
			input: map[string]interface{}{
				"username": "testuser",
				"password": "wrongpass",
			},
			setupMock: func(mock *service.MockService) {
				mock.On("Login", "testuser", "wrongpass").Return(uint64(0), "", httperrors.BadRequestError("Invalid username or password"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid username or password\n",
		},
		{
			name: "failure - internal server error",
			input: map[string]interface{}{
				"username": "testuser",
				"password": "wrongpass",
			},
			setupMock: func(mock *service.MockService) {
				mock.On("Login", "testuser", "wrongpass").Return(uint64(0), "", httperrors.InternalServerError("an error occurred while trying to login", errors.New("internal server error")))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "an error occurred while trying to login: internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(service.MockService)
			tt.setupMock(mockService)
			handler := NewHandler(mockService)

			jsonBytes, _ := json.Marshal(tt.input)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBytes))
			w := httptest.NewRecorder()

			handler.Login(w, req)

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
