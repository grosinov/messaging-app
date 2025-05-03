package controller

import (
	"github.com/challenge/pkg/errors"
	"net/http"

	"github.com/challenge/pkg/helpers"
)

type HealthResponse struct {
	Health string
}

// Check returns the health of the service and DB
func (h Handler) Check(w http.ResponseWriter, r *http.Request) {
	err := h.Service.Health()
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	helpers.RespondJSON(w, HealthResponse{Health: "OK"})
}
