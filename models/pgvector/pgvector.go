package models_pgvector

import (
	"time"

	"gorm.io/gorm"
)

type Gplace struct {
	ID              uint           `gorm:"primaryKey;autoIncrement"`
	Name            string         `gorm:"type:text"`
	Address         string         `gorm:"type:text"`
	GmapID          string         `gorm:"type:text;unique"`
	Description     string         `gorm:"type:text"`
	Latitude        float64        `gorm:"type:decimal"`
	Longitude       float64        `gorm:"type:decimal"`
	Category        []string       `gorm:"type:text[]"`
	AvgRating       float64        `gorm:"type:decimal(3,2)"`
	NumOfReviews    int            `gorm:"type:integer"`
	Price           string         `gorm:"type:text"`
	Hours           interface{}    `gorm:"type:jsonb"`
	Misc            interface{}    `gorm:"type:jsonb"`
	State           string         `gorm:"type:text"`
	RelativeResults []string       `gorm:"type:text[]"`
	URL             string         `gorm:"type:text"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type Greview struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	UserID    string         `gorm:"type:text"`
	Name      string         `gorm:"type:text"`
	Time      int64          `gorm:"type:bigint"`
	Rating    int            `gorm:"type:integer;check:rating >= 1 AND rating <= 5"`
	Text      string         `gorm:"type:text"`
	Pics      interface{}    `gorm:"type:jsonb"`
	Resp      interface{}    `gorm:"type:jsonb"`
	GmapID    string         `gorm:"type:text"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Embedding []float32      `gorm:"type:vector(1536)"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
