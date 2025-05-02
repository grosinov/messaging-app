package controller

import (
	"github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/helpers"
	"net/http"
)

type LoginResponse struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

// Login authenticates a user and returns a token
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	id, token, err := h.Service.Login(username, password)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	helpers.RespondJSON(w, LoginResponse{ID: id, Token: token})
}
