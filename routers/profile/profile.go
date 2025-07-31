package routers_profile

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_profile "pashmak.com/pashmak/controllers/profile"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	middlewares_validation "pashmak.com/pashmak/middlewares/validation"
	serializers_profile "pashmak.com/pashmak/serializers/profile"
	services_openai "pashmak.com/pashmak/services/openai"
	services_place "pashmak.com/pashmak/services/place"
	services_profile "pashmak.com/pashmak/services/profile"

	services_auth "pashmak.com/pashmak/services/auth"
)

func ProfileRoutes(router *gin.Engine, db *gorm.DB, pgvectorDB *gorm.DB, redis *redis.Client, minio *minio.Client, appConfig *bootstrap.AppConfig) {
	profileService := services_profile.NewProfileService(db, minio, appConfig)
	openaiService := services_openai.NewOpenAIService(appConfig.OpenaiApiKey)
	placeService := services_place.NewPlaceService(db, appConfig, openaiService, minio, pgvectorDB)
	profileController := controllers_profile.NewProfileController(profileService, placeService)
	authService := services_auth.NewAuthService(db, redis, appConfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	profile := router.Group("/profiles")
	{
		profile.GET("/me", authMiddleware.LoginMiddleware(), profileController.GetMyProfile)
		profile.POST("/me/update",
			// turn into patch
			middlewares_validation.ValidationMiddleware[serializers_profile.UpdateUserProfileRequest](),
			authMiddleware.LoginMiddleware(), profileController.UpdateUserProfile)
		profile.GET("/:id", profileController.GetProfileByID)
		// saved location apis
		profile.GET("me/saved/label", authMiddleware.LoginMiddleware(), profileController.GetUserPlaceLabels)
		profile.POST("me/saved/label", authMiddleware.LoginMiddleware(), profileController.CreatePlaceLabel)
		profile.DELETE("me/saved/label/:id", authMiddleware.LoginMiddleware(), profileController.DeletePlaceLabel)
		profile.GET("me/saved/location/:place_label_id", authMiddleware.LoginMiddleware(), profileController.GetSavedLocationsByPlaceLabel)
		profile.POST("me/saved/location", authMiddleware.LoginMiddleware(), profileController.CreateSavedLocation)
		profile.PATCH("me/saved/location", authMiddleware.LoginMiddleware(), profileController.UpdateSavedLocation)
		profile.DELETE("me/saved/location/:id", authMiddleware.LoginMiddleware(), profileController.HardDeleteSavedLocation)
		// avatar apis
		profile.GET("/avatar/:file_uuid", profileController.GetUserAvatarObject)
		profile.POST("/avatar/upload", authMiddleware.LoginMiddleware(), profileController.UploadUserAvatar)
		profile.GET("/me/search/history", authMiddleware.LoginMiddleware(), profileController.FetchSearchHistory)
		profile.DELETE("/me/search/history/:id", authMiddleware.LoginMiddleware(), profileController.DeleteSearchHistory)
		profile.DELETE("/me/search/history/clear", authMiddleware.LoginMiddleware(), profileController.ClearSearchHistory)
	}
}
