package routers_comment

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_comment "pashmak.com/pashmak/controllers/comment"
	services_comment "pashmak.com/pashmak/services/comment"
)

func CommentRoutes(router *gin.Engine, db *gorm.DB, appconfig *bootstrap.AppConfig){
	commentService := services_comment.NewCommentService(db, appconfig)
	commentController := controllers_comment.NewCommentController(commentService)

	comment := router.Group("/comments")
	{
		comment.GET("/:token", commentController.GetCommentsByPlaceToken)
	}
}