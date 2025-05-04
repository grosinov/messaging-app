package models

type Message struct {
	Id          uint64  `json:"id" db:"id"`
	SenderId    uint64  `json:"sender_id" gorm:"foreignkey:SenderId" db:"sender_id"`
	RecipientId uint64  `json:"recipient_id" db:"recipient_id"`
	Timestamp   string  `json:"timestamp" db:"timestamp"`
	Content     Content `json:"content" db:"content"`
}

type Content struct {
	Type string `json:"type" db:"type"`
	Text string `json:"text" db:"text"`
}
