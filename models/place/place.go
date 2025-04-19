package models_place

import (
	"time"

	models_comment "pashmak.com/pashmak/models/comment"
)

type Place struct {
	ID        uint                     `gorm:"primaryKey;autoIncrement"`
	// TODO: Add place TOKEN
	Name      string                   `gorm:"not null;size:255"`
	CreatedAt time.Time                `gorm:"autoCreateTime"`
	Comments  []models_comment.Comment `gorm:"foreignKey:PlaceID"` // One-to-many relationship
}
