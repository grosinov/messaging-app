package service

import (
	"github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func (s ServiceImpl) CreateUser(username, password string) (*models.User, error) {
	user, err := s.Repository.GetUserByUsername(username)
	if err == nil {
		return nil, errors.BadRequestError("user already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, errors.InternalServerError("an error occurred while trying to create user", err)
	}

	user, err = s.Repository.CreateUser(models.User{
		Username: username,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, errors.InternalServerError("an error occurred while trying to create user", err)
	}

	return user, err
}

func (s ServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.Repository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
