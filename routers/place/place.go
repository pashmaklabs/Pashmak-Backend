package routers_place

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	"pashmak.com/pashmak/controllers/place"
	"pashmak.com/pashmak/services/comment"
	services_pagination "pashmak.com/pashmak/services/pagination"
	"pashmak.com/pashmak/services/place"
)


func PlaceRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appconfig *bootstrap.AppConfig) {
	paginationService := services_pagination.NewPaginationService(db)
	placeService := services_place.NewPlaceService(db, appconfig)
	commentService := services_comment.NewCommentService(db, appconfig, paginationService)
	placeController := controllers_place.NewPlaceController(placeService, commentService)

	place := router.Group("/places")
	{
		place.GET("/:id", placeController.GetPlace)
		place.GET("/", placeController.SearchPlace)
	}
}