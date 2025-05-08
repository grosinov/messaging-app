package controller

import (
	"encoding/json"
	"errors"
	"github.com/challenge/pkg/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(mockService *service.MockService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success",
			mockSetup: func(mockService *service.MockService) {
				mockService.On("Health").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"health": "ok",
			},
		},
		{
			name: "service error",
			mockSetup: func(mockService *service.MockService) {
				mockService.On("Health").Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error: service error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(service.MockService)
			tt.mockSetup(mockService)

			handler := NewHandler(mockService)
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			handler.Check(w, req)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				response := w.Body.String()
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
