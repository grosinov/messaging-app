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
	requestUser := r.Context().Value("username").(string)
	senderStr := r.FormValue("sender")
	recipientStr := r.FormValue("recipient")
	contentStr := r.FormValue("content")

	senderID, recipientID, content, err := h.validateMessageSending(senderStr, recipientStr, contentStr, requestUser)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	message, err := h.Service.SendMessage(senderID, recipientID, content)
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
	requestUser := r.Context().Value("username").(string)
	recipientStr := r.FormValue("recipient")

	recipient, err := h.validateUserFromStr(recipientStr)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	startStr := r.FormValue("start")
	start, err := strconv.ParseUint(startStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid start value", http.StatusBadRequest)
	}
	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limitStr = "100"
	}
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
	}

	if recipient.Username != requestUser {
		http.Error(w, "You are not allowed to get messages from this user", http.StatusForbidden)
	}

	messages, err := h.Service.GetMessages(recipient.ID, start, limit)

	helpers.RespondJSON(w, messages)
}

func (h Handler) validateMessageSending(senderStr, recipientStr, contentStr, requestUser string) (uint64, uint64, *models.Content, error) {
	sender, err := h.validateUserFromStr(senderStr)
	if err != nil {
		return 0, 0, nil, err
	}

	recipient, err := h.validateUserFromStr(recipientStr)
	if err != nil {
		return 0, 0, nil, err
	}

	var content *models.Content
	err = json.Unmarshal([]byte(contentStr), &content)
	if err != nil {
		return 0, 0, nil, errors.BadRequestError("Invalid content")
	}

	if sender.Username != requestUser {
		return 0, 0, nil, errors.ErrorResponse{
			Status:  http.StatusForbidden,
			Message: "You are not allowed to send messages from this user",
		}
	}

	return sender.ID, recipient.ID, content, nil
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
