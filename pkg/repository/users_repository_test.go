package repository

import (
	"errors"
	"github.com/challenge/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepositoryImpl_CreateUser(t *testing.T) {
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
			user := &models.User{
				Username: "testuser",
				Password: "testpassword",
			}
			user, err := repo.CreateUser(user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)

				user, err = repo.GetUser(user.ID)
				assert.NoError(t, err)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "testpassword", user.Password)
				assert.Equal(t, uint64(1), user.ID)
			}
		})
	}
}

func TestRepositoryImpl_GetUser(t *testing.T) {
	db := SetupTestDB(t)
	repo := &RepositoryImpl{DB: db}
	user, err := repo.CreateUser(&models.User{
		Username: "testuser",
		Password: "testpassword",
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), user.ID)

	tests := []struct {
		name          string
		ID            uint64
		setup         func(t *testing.T) *RepositoryImpl
		expectedError error
	}{
		{
			name: "success",
			ID:   1,
			setup: func(t *testing.T) *RepositoryImpl {
				return repo
			},
			expectedError: nil,
		},
		{
			name: "failure - user not found",
			setup: func(t *testing.T) *RepositoryImpl {
				return repo
			},
			expectedError: errors.New("record not found"),
		},
		{
			name: "failure - database connection is closed",
			ID:   1,
			setup: func(t *testing.T) *RepositoryImpl {
				db = SetupTestDBConnectionClosed(t)
				return &RepositoryImpl{DB: db}
			},
			expectedError: errors.New("sql: database is closed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := tt.setup(t)
			user, err := repository.GetUser(tt.ID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint64(1), user.ID)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "testpassword", user.Password)
			}
		})
	}
}

func TestRepositoryImpl_GetUserByUsername(t *testing.T) {
	db := SetupTestDB(t)
	repo := &RepositoryImpl{DB: db}
	user, err := repo.CreateUser(&models.User{
		Username: "testuser",
		Password: "testpassword",
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), user.ID)

	tests := []struct {
		name          string
		username      string
		setup         func(t *testing.T) *RepositoryImpl
		expectedError error
	}{
		{
			name:     "success",
			username: "testuser",
			setup: func(t *testing.T) *RepositoryImpl {
				return repo
			},
			expectedError: nil,
		},
		{
			name:     "failure - user not found",
			username: "nonexistent",
			setup: func(t *testing.T) *RepositoryImpl {
				return repo
			},
			expectedError: errors.New("record not found"),
		},
		{
			name:     "failure - database connection is closed",
			username: "testuser",
			setup: func(t *testing.T) *RepositoryImpl {
				db = SetupTestDBConnectionClosed(t)
				return &RepositoryImpl{DB: db}
			},
			expectedError: errors.New("sql: database is closed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := tt.setup(t)
			user, err := repository.GetUserByUsername(tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint64(1), user.ID)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "testpassword", user.Password)
			}
		})
	}
}
