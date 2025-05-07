package service

import (
	"errors"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

func TestServiceImpl_CreateUser(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	svc := NewService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		username      string
		password      string
		mockBehavior  func()
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:     "success",
			username: "testuser",
			password: "password123",
			mockBehavior: func() {
				expectedUser := models.User{
					Username: "testuser",
					Password: string(hashedPassword),
				}
				mockRepo.On("GetUserByUsername", "testuser").Return(nil, gorm.ErrRecordNotFound).Once()
				mockRepo.On("CreateUser", mock.Anything).Return(&expectedUser, nil).Once()
			},
			expectedUser: &models.User{
				Username: "testuser",
				Password: string(hashedPassword),
			},
			expectedError: nil,
		},
		{
			name:     "user already exists",
			username: "testuser",
			password: "password123",
			mockBehavior: func() {
				expectedUser := models.User{
					Username: "testuser",
					Password: string(hashedPassword),
				}
				mockRepo.On("GetUserByUsername", "testuser").Return(&expectedUser, nil).Once()
			},
			expectedUser:  nil,
			expectedError: httperrors.BadRequestError("user already exists"),
		},
		{
			name:     "user already exists",
			username: "testuser",
			password: "password123",
			mockBehavior: func() {
				mockRepo.On("GetUserByUsername", "testuser").Return(nil, errors.New("unexpected error")).Once()
			},
			expectedUser:  nil,
			expectedError: httperrors.InternalServerError("an error occurred while trying to create user", errors.New("unexpected error")),
		},
		{
			name:     "failed to create user",
			username: "testuser",
			password: "password123",
			mockBehavior: func() {
				mockRepo.On("GetUserByUsername", "testuser").Return(nil, gorm.ErrRecordNotFound).Once()
				mockRepo.On("CreateUser", mock.Anything).Return(nil, errors.New("database error")).Once()
			},
			expectedUser:  nil,
			expectedError: httperrors.InternalServerError("an error occurred while trying to create user", errors.New("database error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockBehavior()

			user, err := svc.CreateUser(tt.username, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Password, user.Password)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
			}
		})
	}
}

func TestServiceImpl_GetUser(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	svc := NewService(mockRepo)

	tests := []struct {
		name          string
		userID        uint64
		mockSetup     func(m *repository.MockRepository)
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:   "Successful Get User",
			userID: 1,
			mockSetup: func(m *repository.MockRepository) {
				expectedUser := &models.User{ID: 1, Username: "testuser"}
				m.On("GetUser", uint64(1)).Return(expectedUser, nil)
			},
			expectedUser:  &models.User{ID: 1, Username: "testuser"},
			expectedError: nil,
		},
		{
			name:   "User Not Found",
			userID: 1,
			mockSetup: func(m *repository.MockRepository) {
				m.On("GetUser", uint64(1)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
		{
			name:   "failed to get user",
			userID: 1,
			mockSetup: func(m *repository.MockRepository) {
				m.On("GetUser", uint64(1)).Return(nil, errors.New("db error")).Once()
			},
			expectedUser:  nil,
			expectedError: errors.New("an error occurred while trying to login: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup(mockRepo)

			user, err := svc.GetUser(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestServiceImpl_GetUserByUsername(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	svc := NewService(mockRepo)

	tests := []struct {
		name          string
		username      string
		mockSetup     func()
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:     "Successful Get User By Username",
			username: "testuser",
			mockSetup: func() {
				expectedUser := &models.User{Username: "testuser"}
				mockRepo.On("GetUserByUsername", "testuser").Return(expectedUser, nil).Once()
			},
			expectedUser:  &models.User{Username: "testuser"},
			expectedError: nil,
		},
		{
			name:     "User Not Found",
			username: "testuser",
			mockSetup: func() {
				mockRepo.On("GetUserByUsername", "testuser").Return(nil, errors.New("user not found")).Once()
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			user, err := svc.GetUserByUsername(tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}
