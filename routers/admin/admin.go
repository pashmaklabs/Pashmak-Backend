package routers_admin

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
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

func AdminRoutes(router *gin.Engine, db *gorm.DB, pgvectorDB *gorm.DB, redis *redis.Client, minio *minio.Client, appConfig *bootstrap.AppConfig) {
	commentService := services_comment.NewCommentService(db, pgvectorDB, appConfig)
	commentController := controllers_comment.NewCommentController(commentService)
	authService := services_auth.NewAuthService(db, redis, appConfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	admin := router.Group("/admin")
	{
		admin.GET("/reported-comments",
			authMiddleware.LoginMiddleware(),
			authMiddleware.PermissionMiddleware(db, "view_reports"), commentController.GetReportedComments)
		admin.POST("/reported-comments/:id",
			authMiddleware.LoginMiddleware(),
			authMiddleware.PermissionMiddleware(db, "view_reports"),
			middlewares_validation.ValidationMiddleware[serializers_comment.ChangeReportStatus](), commentController.ChangeReportStatus)
	}
}
