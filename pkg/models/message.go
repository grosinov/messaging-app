package models

type Message struct {
	Id          int    `json:"id" db:"id"`
	SenderId    int    `json:"sender_id" db:"sender_id"`
	RecipientId int    `json:"recipient_id" db:"recipient_id"`
	Type        string `json:"type" db:"type"`
	Text        string `json:"text" db:"text"`
}
