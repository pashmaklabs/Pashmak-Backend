package routers_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	controlllers_auth "pashmak.com/pashmak/controllers"
	services_auth "pashmak.com/pashmak/services"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client) {
	routeService := services_auth.NewAuthService(db, redis)
	routeController := controlllers_auth.NewAuthController(routeService)
	
	auth := router.Group("/auth")
	{
		auth.POST("/send-otp", routeController.SendOTP)
	}
}