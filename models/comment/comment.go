package models_comment

import (
	"time"

	models_auth "pashmak.com/pashmak/models/auth"
)

type Comment struct {
	ID      uint             `gorm:"primaryKey;autoIncrement"`
	Content string           `gorm:"not null;size:1000"`
	Rating  uint             `gorm:"default:0"`
	UserID  uint             `gorm:"not null"`
	User    models_auth.User `gorm:"foreignKey:UserID"`
	// Place               // Place Model
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
