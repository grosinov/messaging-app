package repository

import (
	"github.com/challenge/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user models.User) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateLoginSession(login *models.Login) (*models.Login, error)
}

type RepositoryImpl struct {
	DB *gorm.DB
}
