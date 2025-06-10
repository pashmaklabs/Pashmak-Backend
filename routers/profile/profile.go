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
	services_profile "pashmak.com/pashmak/services/profile"

	services_auth "pashmak.com/pashmak/services/auth"
)

func ProfileRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, minio *minio.Client, appConfig *bootstrap.AppConfig) {
	profileService := services_profile.NewProfileService(db, minio, appConfig)
	profileController := controllers_profile.NewProfileController(profileService)
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
		profile.GET("/me/saved-locs", authMiddleware.LoginMiddleware(), profileController.GetSavedLocations)
		profile.POST("/me/saved-locs", authMiddleware.LoginMiddleware(), profileController.AddSavedLocation)
		profile.GET("/avatar/:file_uuid", profileController.GetUserAvatarObject)
		profile.POST("/avatar/upload", authMiddleware.LoginMiddleware(), profileController.UploadUserAvatar)
		profile.GET("/me/search/history", authMiddleware.LoginMiddleware(), profileController.FetchSearchHistory)
	}
}
