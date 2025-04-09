package routers_profile

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_profile "pashmak.com/pashmak/controllers/profile"
	services_profile "pashmak.com/pashmak/services/profile"
)

func ProfileRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appConfig *bootstrap.AppConfig) {
	profileService := services_profile.NewProfileService(db, appConfig)
	profileController := controllers_profile.NewProfilesController(profileService)
	
	profile := router.Group("/profiles")
	{
		profile.GET("/me", routeMiddleware.LoginMiddleware(), profileController.)
	}
}
