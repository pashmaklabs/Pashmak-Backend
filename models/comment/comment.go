package models_comment

import (
	"time"

	models_auth "pashmak.com/pashmak/models/auth"
)

type Comment struct {
	ID        uint             `gorm:"primaryKey;autoIncrement"`
	Content   string          `gorm:"not null;size:1000"`
	Rating    uint           `gorm:"default:0"`
	User      models_auth.User
	// Place               // Place Model
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	Likes     uint          `gorm:"default:0"`
	Dislikes  uint          `gorm:"default:0"`
}
