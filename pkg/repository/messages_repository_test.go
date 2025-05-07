package repository

import (
	"errors"
	"github.com/challenge/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepositoryImpl_SaveMessage(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) *RepositoryImpl
		expectedError error
	}{
		{
			name: "success",
			setup: func(t *testing.T) *RepositoryImpl {
				db := SetupTestDB(t)
				return &RepositoryImpl{DB: db}
			},
			expectedError: nil,
		},
		{
			name: "failure - database connection is closed",
			setup: func(t *testing.T) *RepositoryImpl {
				db := SetupTestDBConnectionClosed(t)
				return &RepositoryImpl{DB: db}
			},
			expectedError: errors.New("sql: database is closed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setup(t)
			timestamp := time.Now().String()
			print(timestamp)
			message := &models.Message{
				SenderID:    1,
				RecipientID: 2,
				Timestamp:   timestamp,
				Content: models.Content{
					Type: "text",
					Text: "test message",
				},
			}
			_, err := repo.SaveMessage(message)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)

				var savedMessage models.Message
				err = repo.DB.First(&savedMessage).Error
				assert.NoError(t, err)
				assert.Equal(t, uint64(1), savedMessage.Id)
				assert.Equal(t, uint64(1), savedMessage.SenderID)
				assert.Equal(t, uint64(2), savedMessage.RecipientID)
				assert.Equal(t, timestamp, savedMessage.Timestamp)
				assert.Equal(t, models.Content{Type: "text", Text: "test message"}, savedMessage.Content)
			}
		})
	}
}

func TestRepositoryImpl_GetMessagesFromUser(t *testing.T) {
	db := SetupTestDB(t)
	repo := &RepositoryImpl{DB: db}
	timestamp := time.Now().String()
	message := &models.Message{
		SenderID:    1,
		RecipientID: 2,
		Timestamp:   timestamp,
		Content: models.Content{
			Type: "text",
			Text: "test message",
		},
	}
	message, err := repo.SaveMessage(message)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), message.Id)

	tests := []struct {
		name          string
		setup         func(t *testing.T) *RepositoryImpl
		recipientID   uint64
		expectedError error
		assertions    func(t *testing.T, messages []models.Message, err error)
	}{
		{
			name: "success",
			setup: func(t *testing.T) *RepositoryImpl {
				return repo
			},
			recipientID:   2,
			expectedError: nil,
			assertions: func(t *testing.T, messages []models.Message, err error) {
				assert.NoError(t, err)
				assert.Len(t, messages, 1)
				assert.Equal(t, uint64(1), messages[0].SenderID)
				assert.Equal(t, uint64(2), messages[0].RecipientID)
				assert.Equal(t, "text", messages[0].Content.Type)
				assert.Equal(t, "test message", messages[0].Content.Text)
			},
		},
		{
			name: "user not found returns empty list",
			setup: func(t *testing.T) *RepositoryImpl {
				return repo
			},
			recipientID:   999,
			expectedError: nil,
			assertions: func(t *testing.T, messages []models.Message, err error) {
				assert.NoError(t, err)
				assert.Len(t, messages, 0)
			},
		},
		{
			name: "failure - database connection is closed",
			setup: func(t *testing.T) *RepositoryImpl {
				db := SetupTestDBConnectionClosed(t)
				return &RepositoryImpl{DB: db}
			},
			recipientID:   2,
			expectedError: errors.New("sql: database is closed"),
			assertions: func(t *testing.T, messages []models.Message, err error) {
				assert.Error(t, err)
				assert.Equal(t, errors.New("sql: database is closed"), err)
				assert.Len(t, messages, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo = tt.setup(t)
			messages, err := repo.GetMessagesFromUser(tt.recipientID, 1, 100)

			tt.assertions(t, messages, err)
		})
	}
}
