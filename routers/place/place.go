package routers_place

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_place "pashmak.com/pashmak/controllers/place"
	services_comment "pashmak.com/pashmak/services/comment"
	services_openai "pashmak.com/pashmak/services/openai"
	services_place "pashmak.com/pashmak/services/place"
)

func PlaceRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appConfig *bootstrap.AppConfig) {
	openaiService := services_openai.NewOpenAIService(appConfig.OpenaiApiKey)
	placeService := services_place.NewPlaceService(db, appConfig, openaiService)
	commentService := services_comment.NewCommentService(db, appConfig)
	placeController := controllers_place.NewPlaceController(placeService, commentService)

	place := router.Group("/places")
	{
		place.GET("/:id", placeController.GetPlace)
		place.GET("/", placeController.SearchPlace)
	}
}
