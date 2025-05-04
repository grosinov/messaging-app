package service

import (
	"errors"
	"github.com/challenge/pkg/auth"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

func TestServiceImpl_Login(t *testing.T) {
	auth.JwtSecret = []byte("testsecret")
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name          string
		username      string
		password      string
		mockSetup     func(username string)
		expectedID    uint64
		expectedError error
	}{
		{
			name:     "successful login",
			username: "testuser",
			password: "password123",
			mockSetup: func(username string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockUser := &models.User{ID: 1, Username: username, Password: string(hashedPassword)}
				mockRepo.On("GetUserByUsername", username).Return(mockUser, nil).Once()
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:     "user not found",
			username: "nonexistent",
			password: "password123",
			mockSetup: func(username string) {
				mockRepo.On("GetUserByUsername", username).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedID:    0,
			expectedError: httperrors.BadRequestError("invalid username or password"),
		},
		{
			name:     "database error",
			username: "testuser",
			password: "password123",
			mockSetup: func(username string) {
				mockRepo.On("GetUserByUsername", username).Return(nil, errors.New("database error")).Once()
			},
			expectedID:    0,
			expectedError: httperrors.InternalServerError("an error occurred while trying to login", errors.New("database error")),
		},
		{
			name:     "wrong password",
			username: "testuser",
			password: "wrongpassword",
			mockSetup: func(username string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				mockUser := &models.User{ID: 1, Username: username, Password: string(hashedPassword)}
				mockRepo.On("GetUserByUsername", username).Return(mockUser, nil).Once()
			},
			expectedID:    0,
			expectedError: httperrors.BadRequestError("invalid username or password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(tt.username)

			userID, token, err := service.Login(tt.username, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				parsedToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
					return auth.JwtSecret, nil
				})

				assert.Equal(t, tt.expectedID, uint64(parsedToken.Claims.(jwt.MapClaims)["user_id"].(float64)))
				assert.NotEmpty(t, parsedToken.Claims.(jwt.MapClaims)["exp"])
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
			assert.Equal(t, tt.expectedID, userID)
		})
	}
}
