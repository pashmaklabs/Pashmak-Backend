package models_auth

import(
  "time"
)

type User struct {
	ID			uint		`gorm:"primaryKey;autoIncrement"`
	FirstName	string		`gorm:"default:''"`
	LastName	string		`gorm:"default:''"`
	Email		string		`gorm:"unique;not null"`
	Password	string		`gorm:"default:''"`
	CreatedAt	time.Time	`gorm:"autoCreateTime"`
	UpdatedAt	time.Time	`gorm:"autoUpdateTime"`
}

type JWTBlacklist struct{
	ID			uint		`gorm:"primaryKey;autoIncrement"`
	JIT			string		`gorm:"unique;not null"`
	ExpiresAt	int64		`gorm:"not null"` 
}