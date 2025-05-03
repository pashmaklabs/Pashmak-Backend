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
	services_pagination "pashmak.com/pashmak/services/pagination"
)

func CommentRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appconfig *bootstrap.AppConfig) {
	paginationService := services_pagination.NewPaginationService(db)
	commentService := services_comment.NewCommentService(db, appconfig, paginationService)
	commentController := controllers_comment.NewCommentController(commentService)
	authService := services_auth.NewAuthService(db, redis, appconfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	comment := router.Group("/comments")
	{
		comment.GET("/:placeToken", commentController.GetCommentsByPlaceToken)
		comment.POST("/:token/reaction/set",
			middlewares_validation.ValidationMiddleware[serializers_comment.AddReactionRequest](),
			authMiddleware.LoginMiddleware(), commentController.SetNewReaction)
		comment.POST("/:token/reaction/remove",
			authMiddleware.LoginMiddleware(), commentController.RemoveReaction)
		comment.POST("/:token/add-comment",
			middlewares_validation.ValidationMiddleware[serializers_comment.AddCommentRequest](),
			authMiddleware.LoginMiddleware(), commentController.AddNewComment)
	}
}
