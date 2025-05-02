package controller

import (
	"net/http"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// Login authenticates a user and returns a token
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: User must login and a token must be generated

	helpers.RespondJSON(w, models.Login{})
}
