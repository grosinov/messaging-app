package service

import (
	"github.com/challenge/pkg/errors"
	"net/http"
)

func (s ServiceImpl) Health() error {
	if err := s.Repository.HealthCheck(); err != nil {
		return errors.ErrorResponse{
			Status:  http.StatusServiceUnavailable,
			Message: "DB is not available: " + err.Error(),
		}
	}

	return nil
}
