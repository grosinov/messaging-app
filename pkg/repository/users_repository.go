package repository

import (
	"github.com/challenge/pkg/models"
)

func (r RepositoryImpl) CreateUser(user models.User) (*models.User, error) {
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r RepositoryImpl) GetUser(id uint64) (*models.User, error) {
	var user models.User
	if err := r.DB.
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r RepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
