package serializers_comment

import (
	"time"
)

// type CommentResponse struct {
// 	ID        uint                    `json:"id"`
// 	Content   string                  `json:"content"`
// 	Rating    uint                    `json:"rating"`
// 	UserID    uint                    `json:"user_id"`
// 	User      models_auth.User        `json:"user"`
// 	PlaceID   uint                    `json:"place_id"`
// 	CreatedAt time.Time               `json:"created_at"`
// 	Reactions []models_place.Reaction `json:"reactions"`
// }

type CommentResponse struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	Rating    uint         `json:"rating"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar_url"`
}


type AddCommentRequest struct{
	Content string	`json:"content"`
	Rating	uint	`json:"rating"`
}