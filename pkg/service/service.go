package service

import (
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/repository"
	"github.com/stretchr/testify/mock"
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

func NewService(repo repository.Repository) Service {
	return &ServiceImpl{
		Repository: repo,
	}
}

// MockRepository is a mock repository for testing purposes
type MockRepository struct {
	mock.Mock
}
