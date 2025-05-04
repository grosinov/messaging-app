package service

import (
	"errors"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func (m *MockRepository) SaveMessage(message *models.Message) (*models.Message, error) {
	args := m.Called(message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockRepository) GetMessagesFromUser(id, start, limit uint64) ([]models.Message, error) {
	args := m.Called(id, start, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Message), args.Error(1)
}

func TestSendMessage(t *testing.T) {
	mockRepo := new(MockRepository)

	tests := []struct {
		name          string
		sender        uint64
		recipient     uint64
		content       *models.Content
		setupMocks    func()
		expectedError error
	}{
		{
			name:      "success",
			sender:    1,
			recipient: 2,
			content: &models.Content{
				Type: "text",
				Text: "test message",
			},
			setupMocks: func() {
				mockRepo.On("SaveMessage", mock.Anything).Return(&models.Message{
					Id:          1,
					SenderId:    1,
					RecipientId: 2,
					Content: models.Content{
						Type: "text",
						Text: "test message",
					},
					Timestamp: time.DateTime,
				}, nil).Once()
			},

			expectedError: nil,
		},
		{
			name:      "invalid message type",
			sender:    1,
			recipient: 2,
			content: &models.Content{
				Type: "invalid",
				Text: "test message",
			},
			setupMocks:    func() {},
			expectedError: httperrors.BadRequestError("invalid message type"),
		},
		{
			name:      "repository error",
			sender:    1,
			recipient: 2,
			content: &models.Content{
				Type: "text",
				Text: "test message",
			},
			setupMocks: func() {
				mockRepo.On("SaveMessage", mock.Anything).Return(nil, errors.New("repository error")).Once()
			},
			expectedError: httperrors.InternalServerError("an error occurred while trying to save message", errors.New("repository error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			service := NewService(mockRepo)
			_, err := service.SendMessage(tt.sender, tt.recipient, tt.content)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetMessages(t *testing.T) {
	mockRepo := new(MockRepository)

	tests := []struct {
		name          string
		userID        uint64
		start         uint64
		limit         uint64
		setupMocks    func()
		expectedError error
	}{
		{
			name:   "success",
			userID: 1,
			start:  0,
			limit:  100,
			setupMocks: func() {
				mockRepo.On("GetMessagesFromUser", uint64(1), uint64(0), uint64(100)).Return([]models.Message{
					{
						Id:          1,
						SenderId:    1,
						RecipientId: 2,
						Content: models.Content{
							Type: "text",
							Text: "test message",
						},
						Timestamp: time.DateTime,
					},
				}, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "repository error",
			userID: 1,
			start:  0,
			limit:  100,
			setupMocks: func() {
				mockRepo.On("GetMessagesFromUser", uint64(1), uint64(0), uint64(100)).Return(nil, errors.New("repository error")).Once()
			},
			expectedError: httperrors.InternalServerError("an error occurred while trying to get messages", errors.New("repository error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			service := NewService(mockRepo)
			_, err := service.GetMessages(tt.userID, tt.start, tt.limit)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
