package routers_comment

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_comment "pashmak.com/pashmak/controllers/comment"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	middlewares_ratelimit "pashmak.com/pashmak/middlewares/ratelimit"
	middlewares_validation "pashmak.com/pashmak/middlewares/validation"
	serializers_comment "pashmak.com/pashmak/serializers/comment"
	services_auth "pashmak.com/pashmak/services/auth"
	services_comment "pashmak.com/pashmak/services/comment"
)

func CommentRoutes(router *gin.Engine, db *gorm.DB, pgvectorDB *gorm.DB, redisClient *redis.Client, appconfig *bootstrap.AppConfig) {
	commentService := services_comment.NewCommentService(db, pgvectorDB, appconfig)
	commentController := controllers_comment.NewCommentController(commentService)
	authService := services_auth.NewAuthService(db, redisClient, appconfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	// Posting a comment: 10 / 10 min per IP — prevent spam
	addCommentLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 10, 10*time.Minute, "comment_add", middlewares_ratelimit.KeyByIP)

	// Reactions: 60 / minute per IP — user can rapidly like/unlike
	reactionLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 60, time.Minute, "comment_reaction", middlewares_ratelimit.KeyByIP)

	// Reports: 10 / 10 min per IP — low ceiling, abuse vector
	reportLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 10, 10*time.Minute, "comment_report", middlewares_ratelimit.KeyByIP)

	comment := router.Group("/comments")
	{
		comment.GET("/:placeToken", authMiddleware.AuthOrAnonMiddleware(), commentController.GetCommentsByPlaceToken)

		comment.POST("/:id/reaction/set",
			reactionLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_comment.AddReactionRequest](),
			authMiddleware.LoginMiddleware(),
			commentController.SetNewReaction)

		comment.POST("/:id/reaction/remove",
			reactionLimiter.Middleware(),
			authMiddleware.LoginMiddleware(),
			commentController.RemoveReaction)

		comment.POST("/:id/add-comment",
			addCommentLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_comment.AddCommentRequest](),
			authMiddleware.LoginMiddleware(),
			commentController.AddNewComment)

		comment.POST("/:id/report",
			reportLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_comment.SendReportRequest](),
			authMiddleware.LoginMiddleware(),
			commentController.ReportComment)

		comment.GET("/me", authMiddleware.LoginMiddleware(), commentController.GetCommentsByUser)
	}
}