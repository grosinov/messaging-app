package controller

import (
	"encoding/json"
	"github.com/challenge/pkg/errors"
	"net/http"
	"strconv"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

type MessageResponse struct {
	ID        uint64 `json:"id"`
	Timestamp string `json:"timestamp"`
}

// SendMessage send a message from one user to another
func (h Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	requestUser := r.Context().Value("user_id")
	var req models.Message
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := h.Service.GetUser(req.SenderID)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	_, err = h.Service.GetUser(req.RecipientID)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	if req.SenderID != requestUser {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	message, err := h.Service.SendMessage(req.SenderID, req.RecipientID, &req.Content)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	helpers.RespondJSON(w, MessageResponse{
		ID:        message.Id,
		Timestamp: message.Timestamp,
	})
}

// GetMessages get the messages from the logged user to a recipient
func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	requestUser := r.Context().Value("user_id")
	recipientStr := r.FormValue("recipient")

	recipientID, err := strconv.ParseUint(recipientStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipient ID", http.StatusBadRequest)
		return
	}

	_, err = h.Service.GetUser(recipientID)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	startStr := r.FormValue("start")
	start, err := strconv.ParseUint(startStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid start value", http.StatusBadRequest)
		return
	}

	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limitStr = helpers.DefaultMessagesLimit
	}

	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}

	if recipientID != requestUser {
		http.Error(w, "You are not allowed to get messages from this user", http.StatusForbidden)
		return
	}

	messages, err := h.Service.GetMessages(recipientID, start, limit)

	helpers.RespondJSON(w, messages)
}

func (h Handler) validateUserFromStr(userIDstr string) (*models.User, error) {
	userID, err := strconv.ParseUint(userIDstr, 10, 64)
	if err != nil {
		return nil, errors.BadRequestError("Invalid user ID")
	}

	user, err := h.Service.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
