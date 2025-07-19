package routers_place

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_place "pashmak.com/pashmak/controllers/place"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	services_auth "pashmak.com/pashmak/services/auth"
	services_comment "pashmak.com/pashmak/services/comment"
	services_openai "pashmak.com/pashmak/services/openai"
	services_place "pashmak.com/pashmak/services/place"
	services_profile "pashmak.com/pashmak/services/profile"
)

func PlaceRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, minio *minio.Client, appConfig *bootstrap.AppConfig) {
	openaiService := services_openai.NewOpenAIService(appConfig.OpenaiApiKey)
	placeService := services_place.NewPlaceService(db, appConfig, openaiService, minio)
	commentService := services_comment.NewCommentService(db, appConfig)
	profileService := services_profile.NewProfileService(db, minio, appConfig)
	placeController := controllers_place.NewPlaceController(placeService, commentService, profileService, appConfig)
	authService := services_auth.NewAuthService(db, redis, appConfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	place := router.Group("/places")
	{
		place.GET("/:id", authMiddleware.AuthOrAnonMiddleware(), placeController.GetPlace)
		place.GET("/", authMiddleware.AuthOrAnonMiddleware(), placeController.SearchPlace)
		place.POST("/:id/images", placeController.UploadPlaceImage)
		place.GET("/:id/images/:image_name", placeController.GetPlaceImage)
		place.POST("/new_place", placeController.AddNewPlace)
	}
}
