package models

type Content struct {
	Type string `json:"type" db:"type"`
	Text string `json:"text" db:"text"`
}
