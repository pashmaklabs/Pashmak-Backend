package services_comment

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
)

type CommentService struct {
	DB        *gorm.DB
	AppConfig *bootstrap.AppConfig
}


func NewCommentService(db *gorm.DB, appconfig *bootstrap.AppConfig) *CommentService{
	return &CommentService{
		DB : db,
		AppConfig: appconfig,
	}
}