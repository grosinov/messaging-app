package service

import (
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/repository"
)

type Service interface {
	Health() error
	CreateUser(username, password string) (*models.User, error)
	GetUser(id uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	Login(username, password string) (uint64, string, error)
	SendMessage(sender, recipient uint64, content *models.Content) (*models.Message, error)
	GetMessages(id, start, limit uint64) ([]models.Message, error)
}

type ServiceImpl struct {
	Repository repository.Repository
}
