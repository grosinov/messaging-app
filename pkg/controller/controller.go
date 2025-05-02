package controller

import "github.com/challenge/pkg/service"

// Handler provides the interface to handle different requests
type Handler struct {
	service service.Service
}
