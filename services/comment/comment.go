package services_comment

import (
	"errors"

	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_place "pashmak.com/pashmak/models/place"
	serializers_comment "pashmak.com/pashmak/serializers/comment"
	// serializers_comment "pashmak.com/pashmak/serializers/comment"
)

type CommentService struct {
	DB        *gorm.DB
	AppConfig *bootstrap.AppConfig
}

func NewCommentService(db *gorm.DB, appconfig *bootstrap.AppConfig) *CommentService {
	return &CommentService{
		DB:        db,
		AppConfig: appconfig,
	}
}

func (cs *CommentService) GetCommentsByPlaceToken(token string) ([]serializers_comment.CommentResponse, error) {
	var comments []models_place.Comment

	err := cs.DB.
		Select("comments.id, comments.content, comments.rating, comments.user_id, comments.place_id, comments.created_at").
		Where("place_id = ?", token).
		Preload("Place").
		Preload("User").
		Preload("Reactions").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, errors.New("no comments found")
	}

	commentDTOs := make([]serializers_comment.CommentResponse, len(comments))
	for i, comment := range comments {
		commentDTOs[i] = serializers_comment.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			Rating:    comment.Rating,
			PlaceID:   comment.PlaceID,
			PlaceName: comment.Place.Name,
			User: serializers_comment.UserResponse{
				ID:        comment.User.ID,
				FirstName: comment.User.FirstName,
				LastName: comment.User.LastName,
                Avatar:    comment.User.Avatar_url,
                
			},
			CreatedAt: comment.CreatedAt,
		}
	}

	return commentDTOs, nil
}
