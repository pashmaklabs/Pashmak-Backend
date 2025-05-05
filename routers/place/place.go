package routers_place

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	"pashmak.com/pashmak/controllers/place"
	"pashmak.com/pashmak/services/comment"
	"pashmak.com/pashmak/services/place"
)


func PlaceRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appconfig *bootstrap.AppConfig) {
	placeService := services_place.NewPlaceService(db, appconfig)
	commentService := services_comment.NewCommentService(db, appconfig)
	placeController := controllers_place.NewPlaceController(placeService, commentService)

	place := router.Group("/places")
	{
		place.GET("/:id", placeController.GetPlace)
		place.GET("/", placeController.SearchPlace)
	}
}