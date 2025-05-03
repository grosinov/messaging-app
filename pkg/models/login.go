package models

import "time"

type Login struct {
	ID        uint64
	UserID    uint64 `gorm:"foreignkey:UserID"`
	TokenID   string
	ExpiresAt time.Time
}
