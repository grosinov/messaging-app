package service

import (
	"errors"
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s ServiceImpl) CreateUser(username, password string) (*models.User, error) {
	user, err := s.Repository.GetUserByUsername(username)
	if err == nil {
		return nil, httperrors.BadRequestError("user already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, httperrors.InternalServerError("an error occurred while trying to create user", err)
	}

	user, err = s.Repository.CreateUser(models.User{
		Username: username,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, httperrors.InternalServerError("an error occurred while trying to create user", err)
	}

	return user, err
}

func (s ServiceImpl) GetUser(id uint64) (*models.User, error) {
	user, err := s.Repository.GetUser(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperrors.BadRequestError("invalid username or password")
	} else if err != nil {
		return nil, httperrors.InternalServerError("an error occurred while trying to login", err)
	}

	return user, nil
}

func (s ServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.Repository.GetUserByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperrors.BadRequestError("invalid username or password")
	} else if err != nil {
		return nil, httperrors.InternalServerError("an error occurred while trying to login", err)
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
