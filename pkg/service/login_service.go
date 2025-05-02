package service

import (
	"github.com/challenge/pkg/auth"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s ServiceImpl) Login(username, password string) (uint, string, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return 0, "", httperrors.InternalServerError("an error occurred while trying to login", err)
	}
	if user == nil {
		return 0, "", httperrors.BadRequestError("invalid username or password")
	}

	if err = checkPassword(user, password); err != nil {
		return 0, "", httperrors.BadRequestError("invalid username or password")
	}

	tokenID := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"jti":     tokenID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(auth.JwtSecret)

	login := &models.Login{
		UserID:    user.ID,
		TokenID:   tokenID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	session, err := s.Repository.CreateLoginSession(login)
	if err != nil {
		return 0, "", httperrors.InternalServerError("an error occurred while trying to login", err)
	}

	return session.ID, signed, nil
}

func checkPassword(user *models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
