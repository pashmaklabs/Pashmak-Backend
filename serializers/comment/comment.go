package serializers_comment

import (
	"time"

	models_auth "pashmak.com/pashmak/models/auth"
	models_comment "pashmak.com/pashmak/models/comment"
)

type CommentResponse struct {
	ID        uint
	Content   string
	Rating    uint
	User      models_auth.User
	PlaceID   uint
	CreatedAt time.Time
	Reactions []models_comment.Reaction
}
