package services_comment

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rosberry/go-pagination"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_auth "pashmak.com/pashmak/models/auth"
	models_place "pashmak.com/pashmak/models/place"
	serializers_comment "pashmak.com/pashmak/serializers/comment"
	services_auth "pashmak.com/pashmak/services/auth"
	services_paginator "pashmak.com/pashmak/services/pagination"
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

func (cs *CommentService) PaginateComments(c *gin.Context, comments *gorm.DB) (*pagination.Paginator, []serializers_comment.CommentResponse, error){
	pagedComments, paginator, err := services_paginator.Paginate[models_place.Comment](comments, c, cs.DB, 20)
	if err != nil {
		return nil, nil, err
	}

	commentDTOs := make([]serializers_comment.CommentResponse, len(pagedComments))
	for i, comment := range pagedComments {
		likes, dislikes, err := cs.FetchReactionsFromDatabase(comment.ID)
		if err != nil{
			return nil, nil, err
		}
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
			Likes: likes,
			Dislikes: dislikes,
		}
	}

	return paginator, commentDTOs, nil
}

func (cs *CommentService) FetchReactionsFromDatabase(commentID uint)(uint, uint, error){
	var likes uint
	var dislikes uint
	query := `
		SELECT 
			COUNT(*)
		FROM reactions
		WHERE reaction_type = ? AND comment_id = ?` 

	if err := cs.DB.Raw(query, 0, commentID).Scan(&likes).Error; err != nil{
		return 0, 0, nil
	}
	if err := cs.DB.Raw(query, 1, commentID).Scan(&dislikes).Error; err != nil{
		return 0, 0, nil
	}
	return likes, dislikes, nil
}

func (cs *CommentService) GetCommentsByPlaceToken(c *gin.Context, token string) (*pagination.Paginator, []serializers_comment.CommentResponse, error) {
	var comments []models_place.Comment
	commentsQuery := cs.DB.
		Where("place_id = ?", token).
		Preload("User").      // use if you want to use comment.User
		Preload("Reactions"). // use if you want to use comment.Reactions
		Find(&comments)

	if commentsQuery.Error != nil {
		return nil, nil, commentsQuery.Error
	}

	if len(comments) == 0 {
		return nil, nil, errors.New("no comments found")
	}

	paginator, commentDTOs, err := cs.PaginateComments(c, commentsQuery)
	if err != nil{
		return nil, nil, err
	}
	return paginator, commentDTOs, nil
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
	// FIXME: AddComment api has body error
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


func (cs *CommentService) AddReaction(userInfo services_auth.UserInfo, commentToken string, reactionType uint) error{
	var comment models_place.Comment
    if err := cs.DB.First(&comment, commentToken).Error; err != nil {
        return errors.New("comment not found")
    }

	var existingReaction models_place.Reaction
    if err := cs.DB.Where("user_id = ? AND comment_id = ?", userInfo.ID, commentToken).
        First(&existingReaction).Error; err == nil {
        // Update existing reaction
        existingReaction.ReactionType = reactionType
        cs.DB.Save(&existingReaction)
        return nil
    }

	newReaction := models_place.Reaction{
		CommentID: commentToken,
		ReactionType: reactionType,
		UserID: userInfo.ID,
	}
	if err := cs.DB.Create(&newReaction).Error; err != nil{
		return err
	}
	
	return nil
}

func (cs *CommentService) RemoveRection(userInfo services_auth.UserInfo, commentToken string) error{
	var comment models_place.Comment
	if err := cs.DB.First(&comment, commentToken).Error; err != nil{
		return errors.New("comment not found")
	}

	var existingReaction models_place.Reaction
	if err := cs.DB.Where("comment_id = ? AND user_id = ?", commentToken, userInfo.ID).Delete(&existingReaction).Error; err != nil{
		return err
	}
	return nil
}