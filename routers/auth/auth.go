package routers_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_auth "pashmak.com/pashmak/controllers/auth"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	services_auth "pashmak.com/pashmak/services/auth"

	
)

func AuthRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appConfig *bootstrap.AppConfig) {
	authService := services_auth.NewAuthService(db, redis, appConfig)
	authController := controllers_auth.NewAuthController(authService, appConfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)
	
	auth := router.Group("/auth")
	{
		auth.POST("/otp/send", authController.SendOTP)
		auth.POST("/otp/verify", authController.VerifyOTP)
		auth.GET("/protected", authMiddleware.LoginMiddleware(), authController.ProtectedRouter)
		auth.POST("/password", authController.LoginWithPassword)
		auth.POST("/password/forget/send", authController.ForgetPassword)
		auth.POST("/password/forget/verify", authController.ForgetPasswordVerify)
		auth.POST("/password/forget/reset", authMiddleware.LoginMiddleware(), authController.ForgetPasswordReset)
		auth.PATCH("/signup", authMiddleware.LoginMiddleware(), authController.SignUp)
	}


}