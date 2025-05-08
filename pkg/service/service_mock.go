package service

import (
	"github.com/challenge/pkg/models"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Health() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockService) CreateUser(username, password string) (*models.User, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) GetUser(id uint64) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) Login(username, password string) (uint64, string, error) {
	args := m.Called(username, password)
	return args.Get(0).(uint64), args.String(1), args.Error(2)
}

func (m *MockService) SendMessage(sender, recipient uint64, content *models.Content) (*models.Message, error) {
	args := m.Called(sender, recipient, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockService) GetMessages(id, start, limit uint64) ([]models.Message, error) {
	args := m.Called(id, start, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Message), args.Error(1)
}
