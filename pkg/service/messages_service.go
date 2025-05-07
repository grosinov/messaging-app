package service

import (
	httperrors "github.com/challenge/pkg/errors"
	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
	"time"
)

func (s ServiceImpl) SendMessage(sender, recipient uint64, content *models.Content) (*models.Message, error) {
	valid := helpers.MessageTypes[content.Type]
	if !valid {
		return nil, httperrors.BadRequestError("invalid message type")
	}

	message := &models.Message{
		SenderID:    sender,
		RecipientID: recipient,
		Content:     *content,
		Timestamp:   time.Now().String(),
	}

	message, err := s.Repository.SaveMessage(message)
	if err != nil {
		return nil, httperrors.InternalServerError("an error occurred while trying to save message", err)
	}

	return message, nil
}

func (s ServiceImpl) GetMessages(id, start, limit uint64) ([]models.Message, error) {
	messages, err := s.Repository.GetMessagesFromUser(id, start, limit)
	if err != nil {
		return nil, httperrors.InternalServerError("an error occurred while trying to get messages", err)
	}

	return messages, nil
}
