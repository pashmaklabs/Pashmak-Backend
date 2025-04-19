package services_comment

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_comment "pashmak.com/pashmak/models/comment"
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

func (cs *CommentService)GetCommentsByToken(token string){
	result := cs.DB.Model(&models_comment.Comment{}).Where(query interface{}, args ...interface{})
}