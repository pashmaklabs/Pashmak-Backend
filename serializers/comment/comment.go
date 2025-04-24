package serializers_comment

import (
	"time"
)

type CommentResponse struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	Rating    uint         `json:"rating" binding:"required,min=0,max=5"`
	User      UserResponse `json:"user"`
	// 	Reactions []models_place.Reaction `json:"reactions"`
	CreatedAt time.Time    `json:"created_at"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar_url"`
}


type AddCommentRequest struct{
	Content string	`json:"content" binding:""`
	Rating	uint	`json:"rating" binding:"required,numeric,min=0,max=5"`
}