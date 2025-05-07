package service

import (
	"errors"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestServiceImpl_Health(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	svc := NewService(mockRepo)

	tests := []struct {
		name          string
		mockSetup     func()
		expectedError *httperrors.ErrorResponse
	}{
		{
			name: "Successful Health Check",
			mockSetup: func() {
				mockRepo.On("HealthCheck").Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Successful Health Check",
			mockSetup: func() {
				mockRepo.On("HealthCheck").Return(errors.New("mock error")).Once()
			},
			expectedError: &httperrors.ErrorResponse{
				Status:  http.StatusServiceUnavailable,
				Message: "DB is not available: mock error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			err := svc.Health()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Status, http.StatusServiceUnavailable)
				assert.Equal(t, tt.expectedError.Message, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
