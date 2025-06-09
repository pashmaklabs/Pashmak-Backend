package models_place

import (
	"time"

	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models/auth"
)

type Place struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	// TODO: Add place TOKEN
	OsmID     uint    `gorm:"column:osm_id"`
	Name      string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Comments  []Comment `gorm:"foreignKey:PlaceID"` // INFO: One-to-many relationship
	Images    []Image   `gorm:"foreignKey:PlaceID"`
}

type Comment struct {
	ID        uint             `gorm:"primaryKey;autoIncrement"`
	Content   string           `gorm:"not null;size:1000"`
	Rating    uint             `gorm:"default:0"`
	UserID    uint             `gorm:"not null"`
	User      models_auth.User `gorm:"foreignKey:UserID"`
	PlaceID   uint             `gorm:"not null"`
	Place     Place            `gorm:"foreignKey:PlaceID"`
	Reactions []Reaction       `gorm:"foreignKey:CommentID"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
}

type Reaction struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CommentID    uint      `gorm:"not null"`
	ReactionType uint      `gorm:"not null"`
	UserID       uint      `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type Image struct {
	gorm.Model
	PlaceID uint
	URL     string `gorm:"not null"`
	AltText string `gorm:"size:100"`
	Caption string `gorm:"size:1000"`
}
