package routers_place

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_place "pashmak.com/pashmak/controllers/place"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	middlewares_ratelimit "pashmak.com/pashmak/middlewares/ratelimit"
	middlewares_validation "pashmak.com/pashmak/middlewares/validation"
	serializers_place "pashmak.com/pashmak/serializers/place"
	services_auth "pashmak.com/pashmak/services/auth"
	services_comment "pashmak.com/pashmak/services/comment"
	services_openai "pashmak.com/pashmak/services/openai"
	services_place "pashmak.com/pashmak/services/place"
	services_profile "pashmak.com/pashmak/services/profile"
)

func PlaceRoutes(router *gin.Engine, db *gorm.DB, pgvectorDB *gorm.DB, redisClient *redis.Client, minioClient *minio.Client, appConfig *bootstrap.AppConfig) {
	openaiService := services_openai.NewOpenAIService(appConfig.OpenaiApiKey)
	placeService := services_place.NewPlaceService(db, appConfig, openaiService, minioClient, pgvectorDB)
	commentService := services_comment.NewCommentService(db, pgvectorDB, appConfig)
	profileService := services_profile.NewProfileService(db, minioClient, appConfig)
	placeController := controllers_place.NewPlaceController(placeService, commentService, profileService, appConfig)
	authService := services_auth.NewAuthService(db, redisClient, appConfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	searchLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 30, time.Minute, "place_search", middlewares_ratelimit.KeyByIP)

	getPlaceLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 120, time.Minute, "place_get", middlewares_ratelimit.KeyByIP)

	imageUploadLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 20, time.Minute, "place_image_upload", middlewares_ratelimit.KeyByIP)

	addPlaceLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 10, 10*time.Minute, "place_add", middlewares_ratelimit.KeyByIP)

	recommendLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 30, time.Minute, "place_recommend", middlewares_ratelimit.KeyByIP)

	place := router.Group("/places")
	{
		place.GET("/:id", getPlaceLimiter.Middleware(), authMiddleware.AuthOrAnonMiddleware(), placeController.GetPlace)
		place.GET("/", searchLimiter.Middleware(), authMiddleware.AuthOrAnonMiddleware(), placeController.SearchPlace)
		place.GET("/recommendations", recommendLimiter.Middleware(), placeController.GetPlaceRecommendations)

		place.POST("/:id/images", imageUploadLimiter.Middleware(), placeController.UploadPlaceImage)
		place.GET("/:id/images/:image_name", getPlaceLimiter.Middleware(), placeController.GetPlaceImage)

		place.POST("/new_place",
			addPlaceLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_place.AddPlaceRequest](),
			authMiddleware.LoginMiddleware(),
			placeController.AddNewPlace)
	}
}