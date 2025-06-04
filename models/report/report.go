package models_report

import (
	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models/auth"
	models_place "pashmak.com/pashmak/models/place"
)

type Report struct {
	gorm.Model
	CommentID uint
	Comment   models_place.Comment
	UserID    uint
	User      models_auth.User
	Reason    string
	Status    string
}
