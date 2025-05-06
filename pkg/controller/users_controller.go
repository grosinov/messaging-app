package controller

import (
	"encoding/json"
	"github.com/challenge/pkg/errors"
	"net/http"

	"github.com/challenge/pkg/helpers"
)

type UserResponse struct {
	ID uint64 `json:"id"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUser creates a new user
func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	user, err := h.Service.CreateUser(req.Username, req.Password)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	helpers.RespondJSON(w, UserResponse{ID: user.ID})
}
