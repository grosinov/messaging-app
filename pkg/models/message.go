package models

type Message struct {
	Id          uint64  `json:"id"`
	SenderID    uint64  `json:"sender_id" db:"sender_id"`
	Sender      User    `json:"-" gorm:"foreignKey:sender_id"`
	RecipientID uint64  `json:"recipient_id" db:"recipient_id"`
	Recipient   User    `json:"-" gorm:"foreignKey:recipient_id"`
	Timestamp   string  `json:"timestamp"`
	Content     Content `json:"content" gorm:"embedded"`
}

type Content struct {
	Type string `json:"type" db:"type"`
	Text string `json:"text" db:"text"`
}
