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
	models_report "pashmak.com/pashmak/models/report"
	serializers_comment "pashmak.com/pashmak/serializers/comment"
	services_auth "pashmak.com/pashmak/services/auth"
	services_paginator "pashmak.com/pashmak/services/pagination"
	"pashmak.com/pashmak/services/placeOsmUtils"
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

func (cs *CommentService) FetchReactionsFromDatabase(commentID uint) (int64, int64, error) {
	var likes int64
	var dislikes int64
	if err := cs.DB.Model(&models_place.Reaction{}).
		Where("reaction_type = ? AND comment_id = ?", 0, commentID).
		Count(&likes).Error; err != nil {
		return 0, 0, nil
	}
	if err := cs.DB.Model(&models_place.Reaction{}).
		Where("reaction_type = ? AND comment_id = ?", 1, commentID).
		Count(&dislikes).Error; err != nil {
		return 0, 0, nil
	}
	return likes, dislikes, nil
}

func (cs *CommentService) CheckIsReactedByCurrentUser(userpayload services_auth.UserInfo, comment models_place.Comment, reactionType uint) (bool, error) {
	var count int64
	err := cs.DB.Model(&models_place.Reaction{}).
		Where("comment_id = ? AND user_id = ? AND reaction_type = ?", comment.ID, userpayload.ID, reactionType).
		Count(&count).Error
	return count > 0, err
}

func (cs *CommentService) PaginateComments(c *gin.Context, comments *gorm.DB, userpayload services_auth.UserInfo, isLoggedIn bool) (*pagination.Paginator, []serializers_comment.CommentResponse, error) {
	pagedComments, paginator, err := services_paginator.Paginate[models_place.Comment](comments, c, cs.DB, 20)
	if err != nil {
		return nil, nil, err
	}

	// Reverse the pagedComments slice
	for i := 0; i < len(pagedComments)/2; i++ {
		j := len(pagedComments) - 1 - i
		pagedComments[i], pagedComments[j] = pagedComments[j], pagedComments[i]
	}

	commentDTOs := make([]serializers_comment.CommentResponse, len(pagedComments))
	for i, comment := range pagedComments {
		likes, dislikes, err := cs.FetchReactionsFromDatabase(comment.ID)
		if err != nil {
			return nil, nil, err
		}
		isLiked := false
		isDisliked := false
		if isLoggedIn {
			isLiked, err = cs.CheckIsReactedByCurrentUser(userpayload, comment, 0)
			if !isLiked {
				isDisliked, err = cs.CheckIsReactedByCurrentUser(userpayload, comment, 1)
			}
		}
		if err != nil {
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
			CreatedAt:               comment.CreatedAt,
			Likes:                   likes,
			Dislikes:                dislikes,
			IsLikedByCurrentUser:    isLiked,
			IsDislikedByCurrentUser: isDisliked,
		}
	}

	return paginator, commentDTOs, nil
}

func (cs *CommentService) PaginateReportedComments(c *gin.Context, comments *gorm.DB) (*pagination.Paginator, []serializers_comment.ReportedCommentsResponse, error) {
	pagedReports, paginator, err := services_paginator.Paginate[models_report.Report](comments, c, cs.DB, 20)
	if err != nil {
		return nil, nil, err
	}

	// Reverse the pagedComments slice
	for i := 0; i < len(pagedReports)/2; i++ {
		j := len(pagedReports) - 1 - i
		pagedReports[i], pagedReports[j] = pagedReports[j], pagedReports[i]
	}

	reportDTOs := make([]serializers_comment.ReportedCommentsResponse, len(pagedReports))
	for i, report := range pagedReports {
		reportDTOs[i] = serializers_comment.ReportedCommentsResponse{
			ID:     report.ID,
			Reason: report.Reason,
			Status: report.Status,
			Comment: serializers_comment.ReportCommentResponse{
				ID:        report.Comment.ID,
				Content:   report.Comment.Content,
				PlaceID:   report.Comment.PlaceID,
				PlaceName: report.Comment.Place.Name,
				User: serializers_comment.UserResponse{
					ID:        report.User.ID,
					FirstName: report.User.FirstName,
					LastName:  report.User.LastName,
					Avatar:    report.User.Avatar_url,
				},
			},
			CreatedAt: report.CreatedAt,
		}
	}

	return paginator, reportDTOs, nil
}

func (cs *CommentService) GetCommentsByPlaceToken(c *gin.Context, token string, userpayload services_auth.UserInfo, isLoggedIn bool) (*pagination.Paginator, []serializers_comment.CommentResponse, error) {
	var comments []models_place.Comment
	commentsQuery := cs.DB.
		Where("place_id = ?", token).
		Preload("User").      // use if you want to use comment.User
		Preload("Reactions"). // use if you want to use comment.Reactions
		Find(&comments)

	if commentsQuery.Error != nil {
		return nil, nil, commentsQuery.Error
	}

	paginator, commentDTOs, err := cs.PaginateComments(c, commentsQuery, userpayload, isLoggedIn)
	if err != nil {
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

	placeTokenInt, err := strconv.ParseUint(placeToken, 10, 64)
	if err != nil {
		return err
	}

	place, err := placeOsmUtils.ImportFromOSM(uint(placeTokenInt), cs.DB)
	if err != nil {
		return err
	}
	result := cs.DB.Create(&models_place.Comment{
		Content:   payload.Content,
		Rating:    payload.Rating,
		UserID:    user.ID,
		User:      userInfo,
		PlaceID:   place.ID,
		Reactions: []models_place.Reaction{},
	})

	return result.Error
	// FIXME: AddComment api has body error
}

func (cs *CommentService) GetAverageRating(placeToken string) (float64, error) {
	// Check if placeToken is an integer (internal place ID)
	if _, err := strconv.ParseUint(placeToken, 10, 64); err != nil {
		// placeToken is not an integer (external place ID like Google Places ID)
		// External places don't have comments in our system, so return 0
		return 0, nil
	}

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

func (cs *CommentService) AddReaction(userInfo services_auth.UserInfo, commentID uint, reactionType uint) error {
	var comment models_place.Comment
	if err := cs.DB.First(&comment, commentID).Error; err != nil {
		return errors.New("comment not found")
	}

	var existingReaction models_place.Reaction
	if err := cs.DB.Where("user_id = ? AND comment_id = ?", userInfo.ID, commentID).
		First(&existingReaction).Error; err == nil {
		// Update existing reaction
		existingReaction.ReactionType = reactionType
		cs.DB.Save(&existingReaction)
		return nil
	}

	newReaction := models_place.Reaction{
		CommentID:    commentID,
		ReactionType: reactionType,
		UserID:       userInfo.ID,
	}
	if err := cs.DB.Create(&newReaction).Error; err != nil {
		return err
	}

	return nil
}

func (cs *CommentService) RemoveRection(userInfo services_auth.UserInfo, commentID uint) error {
	var comment models_place.Comment
	if err := cs.DB.First(&comment, commentID).Error; err != nil {
		return errors.New("comment not found")
	}

	var existingReaction models_place.Reaction
	if err := cs.DB.Where("comment_id = ? AND user_id = ?", commentID, userInfo.ID).Delete(&existingReaction).Error; err != nil {
		return err
	}
	return nil
}

func (cs *CommentService) ReportComment(userInfo services_auth.UserInfo, commentID int, reason string) error {
	var comment models_place.Comment
	if err := cs.DB.First(&comment, commentID).Error; err != nil {
		return errors.New("comment not found")
	}

	newReport := models_report.Report{
		CommentID: uint(commentID),
		UserID:    userInfo.ID,
		Reason:    reason,
		Status:    "Pending",
	}

	if err := cs.DB.Create(&newReport).Error; err != nil {
		return err
	}

	return nil
}

func (cs *CommentService) GetReportedComments(c *gin.Context, status string) (*pagination.Paginator, []serializers_comment.ReportedCommentsResponse, error) {
	var reportedComments []models_report.Report
	query := cs.DB.Model(&models_report.Report{})
	if status != "" {
		result := query.Where("status = ?", status).Preload("User").Preload("Comment").Preload("Comment.Place").Find(&reportedComments)
		if result.Error != nil {
			return nil, nil, result.Error
		}
	}
	result := query.Preload("User").Preload("Comment").Preload("Comment.Place").Find(&reportedComments)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	if len(reportedComments) == 0 {
		return nil, nil, errors.New("no reports found")
	}

	paginator, commentDTOs, err := cs.PaginateReportedComments(c, result)
	return paginator, commentDTOs, err
}

func (cs *CommentService) ChangeReportStatus(status string, reportId string) error {
	var report models_report.Report
	if err := cs.DB.First(&report, reportId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("report not found")
		}
		return err
	}

	res := cs.DB.Model(&report).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (cs *CommentService) GetCommentsByUser(userInfo services_auth.UserInfo) ([]serializers_comment.UserCommentsRespone, error){
	var userComments []models_place.Comment
	query := cs.DB.Model(&models_place.Comment{})
	result := query.Where("user_id = ?", userInfo.ID).Preload("Place").Preload("User").Find(&userComments)
	if result.Error != nil{
		return []serializers_comment.UserCommentsRespone{}, result.Error
	}

	commentDTOs := make([]serializers_comment.UserCommentsRespone, len(userComments))
	for i, comment := range userComments {
		likes, dislikes, err := cs.FetchReactionsFromDatabase(comment.ID)
		if err != nil{
			return []serializers_comment.UserCommentsRespone{}, err
		}
		commentDTOs[i] = serializers_comment.UserCommentsRespone{
			ID:     comment.ID,
			Content: comment.Content,
			Rating: comment.Rating,
			Likes: likes,
			Dislikes: dislikes,
			PlaceID: comment.Place.ID,
			PlaceName: comment.Place.Name,
			CreatedAt: comment.CreatedAt,
		}
	}
	return commentDTOs, nil
}
