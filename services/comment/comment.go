package services_comment

import (
	"errors"
	"strconv"

	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_auth "pashmak.com/pashmak/models/auth"
	models_place "pashmak.com/pashmak/models/place"
	serializers_comment "pashmak.com/pashmak/serializers/comment"
	services_auth "pashmak.com/pashmak/services/auth"
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
		Where("place_id = ?", token).
		Preload("User").      // use if you want to use comment.User
		Preload("Reactions"). // use if you want to use comment.Reaction
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
			ID:      comment.ID,
			Content: comment.Content,
			Rating:  comment.Rating,
			User: serializers_comment.UserResponse{
				ID:        comment.User.ID,
				FirstName: comment.User.FirstName,
				LastName:  comment.User.LastName,
				Avatar:    comment.User.Avatar_url,
			},
			CreatedAt: comment.CreatedAt,
		}
	}

	return commentDTOs, nil
}

func (cs *CommentService) GetUserByGmail(email string) (models_auth.User, error) {
	var user models_auth.User
	result := cs.DB.First(&user, "email = ?", email)
	return user, result.Error
}

func (cs *CommentService) AddNewComment(placeToken string, user services_auth.UserInfo, payload serializers_comment.AddCommentRequest) error {
	userInfo, err := cs.GetUserByGmail(user.Email)
	if err != nil {
		return err
	}

	placeTokenInt, err := strconv.ParseUint(placeToken, 10, 32)
	if err != nil {
		return err
	}

	result := cs.DB.Create(&models_place.Comment{
		Content:   payload.Content,
		Rating:    payload.Rating,
		UserID:    user.ID,
		User:      userInfo,
		PlaceID:   uint(placeTokenInt),
		Reactions: []models_place.Reaction{},
	})

	return result.Error
}

func (cs *CommentService) GetAverageRating(placeToken string) (float64, error) {
	var result serializers_comment.RatingResponse

	err := cs.DB.Model(&models_place.Comment{}).
		Where("place_id = ?", placeToken).
		Select("AVG(rating) as average_rating, COUNT(*) as count").
		Group("place_id").
		Scan(&result).Error

	if err != nil {
		return 0, errors.New("خطایی رخ داده")
	}

	if result.Count == 0 {
		return 0, errors.New("نظری ثبت نشده")
	}

	return result.AverageRating, nil
}
