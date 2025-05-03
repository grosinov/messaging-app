package models

type User struct {
	ID               uint64    `json:"id"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	ReceivedMessages []Message `gorm:"foreignKey:ReceiverID"`
}
