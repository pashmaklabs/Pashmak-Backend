package services_comment

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_comment "pashmak.com/pashmak/models/comment"
	serializers_comment "pashmak.com/pashmak/serializers/comment"
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

func (cs *CommentService)GetCommentsByPlaceToken(token string)([]serializers_comment.CommentResponse, error){
	// result := cs.DB.Model(&models_comment.Comment{}).Where(query interface{}, args ...interface{})
	var comments []serializers_comment.CommentResponse
    err := cs.DB.Model(&models_comment.Comment{}).
        Preload("User").
        Preload("Reactions").
        Where("place_id = ?", token).
        Order("created_at DESC").
        // Limit(limit).
        // Offset(offset).
        Find(&comments).Error
    return comments, err
}