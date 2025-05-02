package repository

import "github.com/challenge/pkg/models"

func (r RepositoryImpl) CreateLoginSession(login *models.Login) (*models.Login, error) {
	if err := r.DB.Create(&login).Error; err != nil {
		return nil, err
	}

	return login, nil
}
