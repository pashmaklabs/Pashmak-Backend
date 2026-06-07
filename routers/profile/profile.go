package routers_profile

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_profile "pashmak.com/pashmak/controllers/profile"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	middlewares_ratelimit "pashmak.com/pashmak/middlewares/ratelimit"
	middlewares_validation "pashmak.com/pashmak/middlewares/validation"
	serializers_profile "pashmak.com/pashmak/serializers/profile"
	services_auth "pashmak.com/pashmak/services/auth"
	services_openai "pashmak.com/pashmak/services/openai"
	services_place "pashmak.com/pashmak/services/place"
	services_profile "pashmak.com/pashmak/services/profile"
)

func ProfileRoutes(router *gin.Engine, db *gorm.DB, pgvectorDB *gorm.DB, redisClient *redis.Client, minioClient *minio.Client, appConfig *bootstrap.AppConfig) {
	profileService := services_profile.NewProfileService(db, minioClient, appConfig)
	openaiService := services_openai.NewOpenAIService(appConfig.OpenaiApiKey)
	placeService := services_place.NewPlaceService(db, appConfig, openaiService, minioClient, pgvectorDB)
	profileController := controllers_profile.NewProfileController(profileService, placeService)
	authService := services_auth.NewAuthService(db, redisClient, appConfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	readLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 60, time.Minute, "profile_read", middlewares_ratelimit.KeyByIP)

	writeLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 20, time.Minute, "profile_write", middlewares_ratelimit.KeyByIP)

	avatarUploadLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 10, time.Minute, "profile_avatar_upload", middlewares_ratelimit.KeyByIP)

	savedWriteLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 30, time.Minute, "profile_saved_write", middlewares_ratelimit.KeyByIP)

	profile := router.Group("/profiles")
	{
		profile.GET("/me", readLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.GetMyProfile)

		profile.POST("/me/update",
			writeLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_profile.UpdateUserProfileRequest](),
			authMiddleware.LoginMiddleware(),
			profileController.UpdateUserProfile)

		profile.GET("/:id", readLimiter.Middleware(), profileController.GetProfileByID)

		// Saved labels
		profile.GET("me/saved/label", readLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.GetUserPlaceLabels)
		profile.POST("me/saved/label", savedWriteLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.CreatePlaceLabel)
		profile.DELETE("me/saved/label/:id", savedWriteLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.DeletePlaceLabel)

		// Saved locations
		profile.GET("me/saved/location/:place_label_id", readLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.GetSavedLocationsByPlaceLabel)
		profile.POST("me/saved/location", savedWriteLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.CreateSavedLocation)
		profile.PATCH("me/saved/location", savedWriteLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.UpdateSavedLocation)
		profile.DELETE("me/saved/location/:id", savedWriteLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.HardDeleteSavedLocation)

		// Avatar
		profile.GET("/avatar/:file_uuid", readLimiter.Middleware(), profileController.GetUserAvatarObject)
		profile.POST("/avatar/upload", avatarUploadLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.UploadUserAvatar)

		// Search history
		profile.GET("/me/search/history", readLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.FetchSearchHistory)
		profile.DELETE("/me/search/history/:id", writeLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.DeleteSearchHistory)
		profile.DELETE("/me/search/history/clear", writeLimiter.Middleware(), authMiddleware.LoginMiddleware(), profileController.ClearSearchHistory)
	}
}