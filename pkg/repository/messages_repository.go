package repository

import "github.com/challenge/pkg/models"

func (r RepositoryImpl) SaveMessage(message *models.Message) (*models.Message, error) {
	if err := r.DB.Create(&message).Error; err != nil {
		return nil, err
	}

	return message, nil
}

func (r RepositoryImpl) GetMessagesFromUser(id, start, limit uint64) ([]models.Message, error) {
	var messages []models.Message
	if err := r.DB.
		Where("recipient_id = ? AND id BETWEEN ? AND ?", id, start, start+limit-1).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
