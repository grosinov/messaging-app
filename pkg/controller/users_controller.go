package controller

import (
	"github.com/challenge/pkg/errors"
	"net/http"

	"github.com/challenge/pkg/helpers"
)

type UserResponse struct {
	ID uint64 `json:"id"`
}

// CreateUser creates a new user
func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.Service.CreateUser(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		errors.HandleError(w, err)
	}

	helpers.RespondJSON(w, UserResponse{ID: user.ID})
}
