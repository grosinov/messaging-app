package models

import "time"

type Login struct {
	ID        uint
	UserID    uint `gorm:"foreignkey:UserID"`
	TokenID   string
	ExpiresAt time.Time
}
