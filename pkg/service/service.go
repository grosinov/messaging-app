package service

import (
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/repository"
)

type Service interface {
	CreateUser(username, password string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	Login(username, password string) (uint, string, error)
}

type ServiceImpl struct {
	Repository repository.Repository
}
