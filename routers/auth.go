package routers_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controlllers_auth "pashmak.com/pashmak/controllers"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	services_auth "pashmak.com/pashmak/services/auth"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appConfig *bootstrap.AppConfig) {
	routeService := services_auth.NewAuthService(db, redis, appConfig)
	routeController := controlllers_auth.NewAuthController(routeService, appConfig)
	routeMiddleware := middlewares_auth.NewAuthMiddleware(routeService)
	
	auth := router.Group("/auth")
	{
		auth.POST("/otp/send", routeController.SendOTP)
		auth.POST("/otp/verify", routeController.VerifyOTP)
		auth.GET("/protected", routeMiddleware.LoginMiddleware(), routeController.ProtectedRouter)
		auth.POST("/password", routeController.LoginWithPassword)
		auth.POST("/password/forget/send", routeController.ForgetPassword)
		auth.POST("/password/forget/verify", routeController.ForgetPasswordVerify)
		auth.POST("/password/forget/reset", routeController.ForgetPasswordReset)
	}
}