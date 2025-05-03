package models

type Message struct {
	Id          uint64  `json:"id" db:"id"`
	SenderId    uint64  `json:"sender_id" gorm:"foreignkey:SenderId" db:"sender_id"`
	RecipientId uint64  `json:"recipient_id" db:"recipient_id"`
	Timestamp   string  `json:"timestamp" db:"timestamp"`
	Content     Content `json:"content" db:"content"`

	Sender   User `gorm:"foreignKey:SenderID"`
	Receiver User `gorm:"foreignKey:ReceiverID"`
}
