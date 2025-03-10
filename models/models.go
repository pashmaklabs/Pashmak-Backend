package authentication

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

type UserOTP struct{
	ID			uint		`gorm:"primaryKey;autoIncrement"`
	Email 		string		`gorm:"unique;not null"`
	OTP			string		`gorm:"not null"`
	CreatedAt	time.Time	`gorm:"autoCreateTime"`
	UpdatedAt 	time.Time	`gorm:"autoUpdateTime"`
}