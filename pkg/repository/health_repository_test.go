package repository

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheck(t *testing.T) {
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
			err := repo.HealthCheck()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
