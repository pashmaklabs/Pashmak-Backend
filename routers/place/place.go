package routers_place

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/services/place"
	"pashmak.com/pashmak/services/comment"
	"pashmak.com/pashmak/controllers/place"
	"pashmak.com/pashmak/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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