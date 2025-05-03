package repository

import (
	"github.com/challenge/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	HealthCheck() error
	CreateUser(user models.User) (*models.User, error)
	GetUser(id uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateLoginSession(login *models.Login) (*models.Login, error)
	SaveMessage(message *models.Message) (*models.Message, error)
	GetMessagesFromUser(id, start, limit uint64) ([]models.Message, error)
}

type RepositoryImpl struct {
	DB *gorm.DB
}
