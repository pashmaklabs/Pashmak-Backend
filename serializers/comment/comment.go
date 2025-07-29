package serializers_comment

import (
	"time"
)

type CommentResponse struct {
	ID                      uint         `json:"id"`
	Content                 string       `json:"content"`
	Rating                  uint         `json:"rating" binding:"required,min=0,max=5"`
	User                    UserResponse `json:"user"`
	Likes                   int64        `json:"likes"`
	Dislikes                int64        `json:"dislikes"`
	IsLikedByCurrentUser    bool         `json:"isLikedByCurrentUser"`
	IsDislikedByCurrentUser bool         `json:"isDislikedByCurrentUser"`
	CreatedAt               time.Time    `json:"created_at"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar_url"`
}

type RatingResponse struct {
	AverageRating float64 `json:"average_rating"`
	Count         int64   `json:"count"`
}
type AddCommentRequest struct {
	Content string `json:"content" binding:"max=1000"`
	Rating  uint   `json:"rating" binding:"required,numeric,min=0,max=5"`
}

type AddReactionRequest struct {
	ReactionType uint `json:"reaction_type"`
	// FIXME: Validate
}

type SendReportRequest struct {
	Reason string `json:"reason" binding:"required,max=500"`
}

type ReportCommentResponse struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	User      UserResponse `json:"user"`
	PlaceID   uint         `json:"place_id"`
	PlaceName string       `jsonn:"place_name"`
}

type ReportedCommentsResponse struct {
	ID        uint                  `json:"id"`
	Reason    string                `json:"reason"`
	Status    string                `json:"status"`
	Comment   ReportCommentResponse `json:"comment"`
	CreatedAt time.Time             `json:"created_at"`
}

type ChangeReportStatus struct {
	Status string `json:"status" binding:"required,oneof=Pending Resolved Dismissed"`
}

type UserCommentsRespone struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Rating    uint      `json:"rating" binding:"required,min=0,max=5"`
	Likes     int64     `json:"likes"`
	Dislikes  int64     `json:"dislikes"`
	PlaceName string    `json:"place_name"`
	PlaceID   uint      `json:"place_id"`
	CreatedAt time.Time `json:"created_at"`
}
