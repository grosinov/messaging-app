package auth

import (
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createTestToken(claims jwt.MapClaims) string {
	JwtSecret = []byte("testsecret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JwtSecret)
	return tokenString
}

func TestValidateUser(t *testing.T) {
	db := repository.SetupTestDB(t)
	middleware := ValidateUser(db)

	testUser := &models.User{
		ID:       1,
		Username: "testuser",
		Password: "password",
	}
	err := db.Create(testUser).Error
	assert.NoError(t, err)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name: "Valid token",
			token: "Bearer " + createTestToken(jwt.MapClaims{
				"user_id": 1,
				"exp":     time.Now().Add(time.Hour).Unix(),
			}),
			expectedStatus: http.StatusOK,
		},
		{
			name: "No Bearer prefix",
			token: createTestToken(jwt.MapClaims{
				"user_id": 1,
				"exp":     time.Now().Add(time.Hour).Unix(),
			}),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Expired token",
			token: "Bearer " + createTestToken(jwt.MapClaims{
				"user_id": 1,
				"exp":     time.Now().Add(-time.Hour).Unix(),
			}),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid user",
			token: "Bearer " + createTestToken(jwt.MapClaims{
				"user_id": 2,
				"exp":     time.Now().Add(time.Hour).Unix(),
			}),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "user_id not a claim",
			token: "Bearer " + createTestToken(jwt.MapClaims{
				"exp": time.Now().Add(time.Hour).Unix(),
			}),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "exp not a claim",
			token: "Bearer " + createTestToken(jwt.MapClaims{
				"user_id": 1,
			}),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rr := httptest.NewRecorder()

			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value("user_id")
				assert.NotNil(t, userID)
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
