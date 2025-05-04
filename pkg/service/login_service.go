package service

import (
	"errors"
	"github.com/challenge/pkg/auth"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func (s ServiceImpl) Login(username, password string) (uint64, string, error) {
	user, err := s.GetUserByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, "", httperrors.BadRequestError("invalid username or password")
	} else if err != nil {
		return 0, "", httperrors.InternalServerError("an error occurred while trying to login", err)
	}

	if err = checkPassword(user, password); err != nil {
		return 0, "", httperrors.BadRequestError("invalid username or password")
	}

	expTime := time.Now().Add(time.Hour * 24)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(auth.JwtSecret)

	return user.ID, signed, nil
}

func checkPassword(user *models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
