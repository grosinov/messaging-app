package controller

import (
	"encoding/json"
	"github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/helpers"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    uint64 `json:"id"`
	Token string `json:"token"`
}

// Login authenticates a user and returns a token
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	id, token, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	helpers.RespondJSON(w, LoginResponse{ID: id, Token: token})
}
