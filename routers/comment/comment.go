package routers_comment

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_comment "pashmak.com/pashmak/controllers/comment"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	middlewares_validation "pashmak.com/pashmak/middlewares/validation"
	serializers_comment "pashmak.com/pashmak/serializers/comment"
	services_auth "pashmak.com/pashmak/services/auth"
	services_comment "pashmak.com/pashmak/services/comment"
)

func CommentRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appconfig *bootstrap.AppConfig) {
	commentService := services_comment.NewCommentService(db, appconfig)
	commentController := controllers_comment.NewCommentController(commentService)
	authService := services_auth.NewAuthService(db, redis, appconfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	comment := router.Group("/comments")
	{
		comment.GET("/:placeToken", commentController.GetCommentsByPlaceToken)
		comment.POST("/:commentToken/reaction",
		middlewares_validation.ValidationMiddleware[serializers_comment.AddReactionRequest](),
			authMiddleware.LoginMiddleware(), commentController.SetNewReaction)
		comment.POST("/:placeToken/add-comment",
			middlewares_validation.ValidationMiddleware[serializers_comment.AddCommentRequest](),
			authMiddleware.LoginMiddleware(), commentController.AddNewComment)
	}
}
