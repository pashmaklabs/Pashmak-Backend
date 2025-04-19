package models_place

import (
	"time"

	models_auth "pashmak.com/pashmak/models/auth"
)

type Place struct {
	ID        uint                     `gorm:"primaryKey;autoIncrement"`
	// TODO: Add place TOKEN
	Name      string                   `gorm:"not null;size:255"`
	CreatedAt time.Time                `gorm:"autoCreateTime"`
	Comments  []Comment `gorm:"foreignKey:PlaceID"` // One-to-many relationship
}



type Comment struct {
	ID        uint               `gorm:"primaryKey;autoIncrement"`
	Content   string             `gorm:"not null;size:1000"`
	Rating    uint               `gorm:"default:0"`
	UserID    uint               `gorm:"not null"`
	User      models_auth.User   `gorm:"foreignKey:UserID"`
	PlaceID   uint               `gorm:"not null"` // Foreign key for Place
	Place     Place `gorm:"foreignKey:PlaceID"`
	CreatedAt time.Time          `gorm:"autoCreateTime"`
	Reactions []Reaction         `gorm:"foreignKey:CommentID"`
}

type Reaction struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CommentID    uint      `gorm:"not null"`
	ReactionType int       `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}