package controller

import (
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
		expectedBody   string
	}{
		{
			name: "success",
			mockSetup: func(mockService *service.MockService) {
				mockService.On("Health").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"ok"}`,
		},
		{
			name: "service error",
			mockSetup: func(mockService *service.MockService) {
				mockService.On("Health").Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"service error"}`,
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

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
